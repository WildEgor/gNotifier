package middleware

import (
	"net/http"
	"strings"

	"github.com/WildEgor/gNotifier/internal/adapters"
	"github.com/gofiber/fiber/v2"
)

type HealthCheckConfig struct {
	Endpoint           string
	HealthCheckAdapter *adapters.HealthCheckAdapter
}

// Example how we can use middleware
func HealthCheck(cfg *HealthCheckConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		methodName := ctx.Method()
		path := ctx.Path()

		if (methodName == "GET" || methodName == "HEAD") && strings.EqualFold(path, cfg.Endpoint) {
			if cfg.HealthCheckAdapter != nil {
				info := cfg.HealthCheckAdapter.Measure(ctx.Context())

				if info.Status != adapters.StatusOK {
					return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"isOk": false,
						"data": fiber.Map{
							"errors": info.Failures,
						},
					})
				}

				return ctx.Status(http.StatusOK).JSON(fiber.Map{
					"isOk": true,
					"data": fiber.Map{
						"state": info,
					},
				})
			}

			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"isOk": true,
				"data": fiber.Map{},
			})
		}

		return ctx.Next()
	}
}
