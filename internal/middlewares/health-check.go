package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckConfig struct {
	Endpoint string
}

// Example how we can use middleware
func HealthCheck(cfg *HealthCheckConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		methodName := ctx.Method()
		path := ctx.Path()

		if (methodName == "GET" || methodName == "HEAD") && strings.EqualFold(path, cfg.Endpoint) {
			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"isOk": true,
				"data": fiber.Map{},
			})
		}

		return ctx.Next()
	}
}
