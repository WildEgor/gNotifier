package routers

import (
	"fmt"
	"time"

	"github.com/WildEgor/gNotifier/internal/adapters"
	"github.com/WildEgor/gNotifier/internal/config"
	handlers "github.com/WildEgor/gNotifier/internal/handlers/amqp"
	"github.com/WildEgor/gNotifier/internal/services/checks"
	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
)

type AMQPRouter struct {
	notifierHandler   *handlers.NotifierHandler
	amqpConfig        *config.AMQPConfig
	healtCheckAdapter *adapters.HealthCheckAdapter
	amqpAdapter       adapters.IRabbitMQAdapter
}

func NewAMQPRouter(
	notifierHandler *handlers.NotifierHandler,
	amqpConfig *config.AMQPConfig,
	healtCheckAdapter *adapters.HealthCheckAdapter,
	amqpAdapter adapters.IRabbitMQAdapter,
) *AMQPRouter {
	return &AMQPRouter{
		notifierHandler:   notifierHandler,
		amqpConfig:        amqpConfig,
		healtCheckAdapter: healtCheckAdapter,
		amqpAdapter:       amqpAdapter,
	}
}

func (r *AMQPRouter) SetupRoutes() error {
	r.healtCheckAdapter.Register(adapters.HealthConfig{
		Name:      "amqp-health-check",
		Timeout:   time.Duration(10 * time.Second),
		SkipOnErr: false,
		Check: checks.NewRabbitMQCheck(&checks.RabbitMQCheckConfig{
			URI:            r.amqpConfig.URI,
			DialTimeout:    time.Duration(5 * time.Second),
			Exchange:       "health-check",
			Queue:          "health-queue",
			RoutingKey:     "health-check",
			ConsumeTimeout: time.Duration(5 * time.Second),
		}),
	})

	connection, err := rabbitmq.NewConn(
		r.amqpConfig.URI,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed to connect: ", err)
	}
	defer connection.Close()

	consumer, err := rabbitmq.NewConsumer(
		connection,
		func(d rabbitmq.Delivery) (action rabbitmq.Action) {
			fmt.Printf("[AMQPRouter] consumed: %v", string(d.Body))
			return rabbitmq.Ack
		},
		"notifier-queue",
		rabbitmq.WithConsumerOptionsRoutingKey("notifier.send-notification"),
		rabbitmq.WithConsumerOptionsExchangeName("notifications"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed consume: ", err)
	}
	defer consumer.Close()

	return nil
}
