package adapters

import (
	"log"

	"github.com/WildEgor/gNotifier/internal/config"
	"github.com/streadway/amqp"
)

type IRabbitMQAdapter interface {
	GetChannel() *amqp.Channel
}

type RabbitMQAdapter struct {
	amqpConfig *config.AMQPConfig
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQAdapter(
	amqpConfig *config.AMQPConfig,
) *RabbitMQAdapter {
	connection, err := amqp.Dial(amqpConfig.URI)
	if err != nil {
		log.Fatal("[RabbitMQAdapter] Failed to connect: ", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Fatal("[RabbitMQAdapter] Failed to open channel: ", err)
	}

	return &RabbitMQAdapter{
		Connection: connection,
		Channel:    channel,
	}
}

func (r *RabbitMQAdapter) GetChannel() *amqp.Channel {
	return r.Channel
}
