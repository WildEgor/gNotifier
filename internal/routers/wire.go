package routers

import (
	handlers "github.com/WildEgor/gNotifier/internal/handlers"
	"github.com/google/wire"
)

var RoutersSet = wire.NewSet(
	handlers.HandlersSet,
	NewHTTPRouter,
	NewAMQPRouter,
)
