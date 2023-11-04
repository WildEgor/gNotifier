package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/WildEgor/gNotifier/internal/adapters"
	"github.com/gofiber/fiber/v2"
)

type HealthCheckConfig struct {
	Endpoint           string
	HealthCheckAdapter *adapters.HealthCheckAdapter
}

// HealthCheck Example how we can use middleware
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
						"data": info.Failures,
					})
				}

				return ctx.Status(http.StatusOK).JSON(fiber.Map{
					"isOk": true,
					"data": fiber.Map{
						"status":    info.Status,
						"timestamp": info.Timestamp,
						"message":   info.ComponentInfo.Name,
					},
				})
			}

			return ctx.Status(http.StatusOK).JSON(fiber.Map{
				"isOk": true,
				"data": fiber.Map{
					"status":    "OK",
					"timestamp": time.Now(),
					"message":   "",
				},
			})
		}

		return ctx.Next()
	}
}
