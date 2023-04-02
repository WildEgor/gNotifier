package checks

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

const (
	defaultExchange       = "health_check"
	defaultDialTimeout    = 200 * time.Millisecond
	defaultConsumeTimeout = time.Second * 3
)

type RabbitMQCheckConfig struct {
	URI            string
	Exchange       string
	RoutingKey     string
	Queue          string
	ConsumeTimeout time.Duration
	DialTimeout    time.Duration
}

func (c *RabbitMQCheckConfig) initDefaultConfig() {
	if c.Exchange == "" {
		c.Exchange = defaultExchange
	}

	if c.RoutingKey == "" {
		host, err := os.Hostname()
		if nil != err {
			c.RoutingKey = "-unknown-"
		}
		c.RoutingKey = host
	}

	if c.Queue == "" {
		c.Queue = fmt.Sprintf("%s.%s", c.Exchange, c.RoutingKey)
	}

	if c.ConsumeTimeout == 0 {
		c.ConsumeTimeout = defaultConsumeTimeout
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = defaultDialTimeout
	}
}

// NewRabbitMQCheck creates new RabbitMQ health check that verifies the following:
// - connection establishing
// - getting channel from the connection
// - declaring topic exchange
// - declaring queue
// - binding a queue to the exchange with the defined routing key
// - publishing a message to the exchange with the defined routing key
// - consuming published message
func NewRabbitMQCheck(cfg *RabbitMQCheckConfig) func(ctx context.Context) error {
	cfg.initDefaultConfig()

	return func(ctx context.Context) (checkErr error) {
		conn, err := rabbitmq.NewConn(
			cfg.URI,
			rabbitmq.WithConnectionOptionsLogging,
		)
		if err != nil {
			checkErr = fmt.Errorf("[RabbitMQCheck] Failed on connection: %w\n", err)
			return
		}
		defer conn.Close()

		// HINT: Create publisher and consumer
		publisher, err := rabbitmq.NewPublisher(
			conn,
			rabbitmq.WithPublisherOptionsLogging,
			rabbitmq.WithPublisherOptionsExchangeName(cfg.Exchange),
			rabbitmq.WithPublisherOptionsExchangeDeclare,
			rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		)
		if err != nil {
			checkErr = fmt.Errorf("[RabbitMQCheck] Failed on getting channel phase: %w\n", err)
		}
		defer publisher.Close()

		consumer, err := rabbitmq.NewConsumer(
			conn,
			func(d rabbitmq.Delivery) (action rabbitmq.Action) {
				fmt.Printf("[RabbitMQCheck] consumed: %v\n", string(d.Body))
				return rabbitmq.Ack
			},
			cfg.Queue,
			rabbitmq.WithConsumerOptionsRoutingKey(cfg.RoutingKey),
			rabbitmq.WithConsumerOptionsExchangeName(cfg.Exchange),
			rabbitmq.WithConsumerOptionsExchangeDeclare,
			rabbitmq.WithConsumerOptionsExchangeKind("topic"),
			rabbitmq.WithConsumerOptionsBinding(rabbitmq.Binding{
				RoutingKey: cfg.RoutingKey,
				BindingOptions: rabbitmq.BindingOptions{
					Declare: true,
				},
			}),
		)
		if err != nil {
			checkErr = fmt.Errorf("[RabbitMQCheck] Failed consume: %w\n", err)
		}
		defer consumer.Close()

		err = publisher.PublishWithContext(
			context.Background(),
			[]byte("test"),
			[]string{cfg.Exchange},
			rabbitmq.WithPublishOptionsContentType("application/json"),
			rabbitmq.WithPublishOptionsExchange(cfg.Exchange),
		)
		if err != nil {
			checkErr = fmt.Errorf("[RabbitMQCheck] failed publish: %w\n", err)
		}

		return checkErr
	}
}
