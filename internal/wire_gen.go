// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/WildEgor/gNotifier/internal/adapters"
	"github.com/WildEgor/gNotifier/internal/config"
	"github.com/WildEgor/gNotifier/internal/handlers/amqp"
	"github.com/WildEgor/gNotifier/internal/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

// Injectors from server.go:

func NewServer() (*fiber.App, error) {
	appConfig := config.NewAppConfig()
	healthCheckAdapter, err := adapters.NewHealthCheckAdapter()
	if err != nil {
		return nil, err
	}
	httpRouter := routers.NewHTTPRouter(healthCheckAdapter)
	notifierHandler := handlers.NewNotifierHandler()
	amqpConfig := config.NewAMQPConfig()
	amqpRouter := routers.NewAMQPRouter(notifierHandler, amqpConfig, healthCheckAdapter)
	app := NewApp(appConfig, httpRouter, amqpRouter)
	return app, nil
}

// server.go:

var ServerSet = wire.NewSet(AppSet)
