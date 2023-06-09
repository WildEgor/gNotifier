package routers

import (
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
}

func NewAMQPRouter(
	notifierHandler *handlers.NotifierHandler,
	amqpConfig *config.AMQPConfig,
	healtCheckAdapter *adapters.HealthCheckAdapter,
) *AMQPRouter {
	return &AMQPRouter{
		notifierHandler:   notifierHandler,
		amqpConfig:        amqpConfig,
		healtCheckAdapter: healtCheckAdapter,
	}
}

func (r *AMQPRouter) SetupRoutes() error {
	// Add health-check to global HealthCheckAdapter
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

	var err error
	connection, err := rabbitmq.NewConn(
		r.amqpConfig.URI,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed to connect: ", err)
	}
	defer connection.Close()

	notifyConsumer, err := rabbitmq.NewConsumer(
		connection,
		r.notifierHandler.Handle,
		"notifier-queue",
		rabbitmq.WithConsumerOptionsRoutingKey("notifier.send-notification"),
		rabbitmq.WithConsumerOptionsExchangeName("notifications"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed notifyConsumer: ", err)
	}
	defer notifyConsumer.Close()

	return nil
}
