package adapters

import (
	"log"

	"github.com/WildEgor/gNotifier/internal/config"
	"github.com/wagslane/go-rabbitmq"
)

type IRabbitMQAdapter interface {
	GetConnection() *rabbitmq.Conn
	SetConsumer(c *rabbitmq.Consumer)
}

type RabbitMQAdapter struct {
	amqpConfig *config.AMQPConfig
	connection *rabbitmq.Conn
	consumers  *map[string]rabbitmq.Consumer
}

func NewRabbitMQAdapter(
	amqpConfig *config.AMQPConfig,
) *RabbitMQAdapter {
	connection, err := rabbitmq.NewConn(
		amqpConfig.URI,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal("[AMQPRouter] Failed to connect: ", err)
	}
	defer connection.Close()

	return &RabbitMQAdapter{
		connection: connection,
	}
}

func (r *RabbitMQAdapter) GetConnection() *rabbitmq.Conn {
	return r.connection
}

func (r *RabbitMQAdapter) SetConsumer(c *rabbitmq.Consumer) {
	defer c.Close()
}
