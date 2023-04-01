package routers

import (
	"context"
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

func (r *AMQPRouter) SetupRoutes(ctx context.Context) error {
	r.healtCheckAdapter.Register(adapters.HealthConfig{
		Name:      "amqp-health-check",
		Timeout:   time.Duration(5 * time.Second),
		SkipOnErr: false,
		Check: checks.NewRabbitMQCheck(&checks.RabbitMQCheckConfig{
			URI: r.amqpConfig.URI,
		}),
	})

	conn, err := rabbitmq.NewConn(
		r.amqpConfig.URI,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(
		conn,
		r.notifierHandler.Handle,
		"notifier",
		rabbitmq.WithConsumerOptionsRoutingKey("notifier.send-notification"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	return nil
}