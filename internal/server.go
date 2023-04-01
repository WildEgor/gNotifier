//go:build wireinject
// +build wireinject

package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var ServerSet = wire.NewSet(AppSet)

func NewServer() (*fiber.App, error) {
	wire.Build(ServerSet)
	return nil, nil
}
