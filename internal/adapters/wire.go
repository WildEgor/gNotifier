package adapters

import (
	"github.com/google/wire"
)

var AdaptersSet = wire.NewSet(
	NewHealthCheckAdapter,
	NewRabbitMQAdapter,
	wire.Bind(new(IRabbitMQAdapter), new(*RabbitMQAdapter)),
)
