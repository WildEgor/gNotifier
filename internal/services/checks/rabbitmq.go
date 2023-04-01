package checks

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/streadway/amqp"
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
		conn, err := amqp.DialConfig(cfg.URI, amqp.Config{
			Dial: amqp.DefaultDial(cfg.DialTimeout),
		})
		if err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed on dial phase: %w", err)
			return
		}
		defer func() {
			// override checkErr only if there were no other errors
			if err := conn.Close(); err != nil && checkErr == nil {
				checkErr = fmt.Errorf("RabbitMQ health check failed to close connection: %w", err)
			}
		}()

		ch, err := conn.Channel()
		if err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed on getting channel phase: %w", err)
			return
		}
		defer func() {
			// override checkErr only if there were no other errors
			if err := ch.Close(); err != nil && checkErr == nil {
				checkErr = fmt.Errorf("RabbitMQ health check failed to close channel: %w", err)
			}
		}()

		if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed during declaring exchange: %w", err)
			return
		}

		if _, err := ch.QueueDeclare(cfg.Queue, false, false, false, false, nil); err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed during declaring queue: %w", err)
			return
		}

		if err := ch.QueueBind(cfg.Queue, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed during binding: %w", err)
			return
		}

		messages, err := ch.Consume(cfg.Queue, "", true, false, false, false, nil)
		if err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed during consuming: %w", err)
			return
		}

		done := make(chan struct{})

		// HINT: consume messages in goroutine
		go func() {
			// block until: a message is received, or message channel is closed (consume timeout)
			<-messages
			// release the channel resources, and unblock the receive on done below
			close(done)
			// now drain any incidental remaining messages
			for range messages {
			}
		}()

		// HINT: send check messages to the same exchange
		p := amqp.Publishing{Body: []byte(time.Now().Format(time.RFC3339Nano))}
		if err := ch.Publish(cfg.Exchange, cfg.RoutingKey, false, false, p); err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed during publishing: %w", err)
			return
		}

		for {
			select {
			case <-time.After(cfg.ConsumeTimeout):
				checkErr = fmt.Errorf("RabbitMQ health check failed due to consume timeout: %w", err)
				return
			case <-ctx.Done():
				checkErr = fmt.Errorf("RabbitMQ health check failed due "+
					"to health check listener disconnect: %w", ctx.Err())
				return
			case <-done:
				return
			}
		}
	}
}
