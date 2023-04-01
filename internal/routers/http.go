package routers

import (
	middleware "github.com/WildEgor/gNotifier/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type HTTPRouter struct {
}

func NewHTTPRouter() *HTTPRouter {
	return &HTTPRouter{}
}

func (r *HTTPRouter) SetupRoutes(app *fiber.App) error {
	hCfg := middleware.HealthCheckConfig{
		Endpoint: "/api/v1/health/check",
	}

	app.Use(middleware.HealthCheck(&hCfg))

	v1 := app.Group("/api/v1")
	healthCheckController := v1.Group("/health")
	healthCheckController.Get("/check", middleware.HealthCheck(&hCfg))

	return nil
}
