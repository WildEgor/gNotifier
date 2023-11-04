package routers

import (
	"github.com/WildEgor/gNotifier/internal/configs"
	"time"

	"github.com/WildEgor/gNotifier/internal/adapters"
	handlers "github.com/WildEgor/gNotifier/internal/handlers/amqp"
	"github.com/WildEgor/gNotifier/internal/services/checks"
	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
)

type AMQPRouter struct {
	notifierHandler    *handlers.NotifierHandler
	amqpConfig         *configs.AMQPConfig
	healthCheckAdapter *adapters.HealthCheckAdapter

	connection       *rabbitmq.Conn
	notifierConsumer *rabbitmq.Consumer
}

func NewAMQPRouter(
	notifierHandler *handlers.NotifierHandler,
	amqpConfig *configs.AMQPConfig,
	healthCheckAdapter *adapters.HealthCheckAdapter,
) *AMQPRouter {
	return &AMQPRouter{
		notifierHandler:    notifierHandler,
		amqpConfig:         amqpConfig,
		healthCheckAdapter: healthCheckAdapter,
	}
}

func (r *AMQPRouter) SetupRoutes() error {
	// Add health-check to global HealthCheckAdapter
	er := r.healthCheckAdapter.Register(adapters.HealthConfig{
		Name:      "amqp-health-check",
		Timeout:   10 * time.Second,
		SkipOnErr: false,
		Check: checks.NewRabbitMQCheck(&checks.RabbitMQCheckConfig{
			URI:            r.amqpConfig.URI,
			DialTimeout:    5 * time.Second,
			Exchange:       "health-check",
			Queue:          "health-queue",
			RoutingKey:     "health-check",
			ConsumeTimeout: 5 * time.Second,
		}),
	})
	if er != nil {
		return er
	}

	var connection, err = rabbitmq.NewConn(
		r.amqpConfig.URI,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed to connect: ", err)
	}
	r.connection = connection

	notifierConsumer, err := rabbitmq.NewConsumer(
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

	r.notifierConsumer = notifierConsumer

	return nil
}

func (r *AMQPRouter) Close() {
	defer r.connection.Close()
	defer r.notifierConsumer.Close()
}
