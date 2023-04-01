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
	amqpAdapter       adapters.IRabbitMQAdapter
	consumers         map[string]rabbitmq.Consumer
	connection        *rabbitmq.Conn
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
		consumers:         make(map[string]rabbitmq.Consumer),
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

	var err error
	r.connection, err = rabbitmq.NewConn(
		r.amqpConfig.URI,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed to connect: ", err)
	}

	consumer, err := rabbitmq.NewConsumer(
		r.connection,
		r.notifierHandler.Handle,
		"notifier-queue",
		rabbitmq.WithConsumerOptionsRoutingKey("notifier.send-notification"),
		rabbitmq.WithConsumerOptionsExchangeName("notifications"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed consume: ", err)
	}
	r.consumers["notifier-queue::notifier.send-notification::notifications"] = *consumer

	return nil
}

func (r *AMQPRouter) Close() {
	log.Info("[AMQPRouter] Close connection...")
	defer r.connection.Close()
	for _, v := range r.consumers {
		defer v.Close()
	}
}
