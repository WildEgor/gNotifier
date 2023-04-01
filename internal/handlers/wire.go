package handlers

import (
	amqp_handlers "github.com/WildEgor/gNotifier/internal/handlers/amqp"
	http_handlers "github.com/WildEgor/gNotifier/internal/handlers/http"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	amqp_handlers.NewNotifierHandler,
	http_handlers.NewHealthCheckHandler,
)
