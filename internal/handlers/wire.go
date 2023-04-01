package handlers

import (
	amqp_handlers "github.com/WildEgor/gNotifier/internal/handlers/amqp"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	amqp_handlers.NewNotifierHandler,
)
