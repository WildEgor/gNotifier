package routers

import (
	"github.com/WildEgor/gNotifier/internal/adapters"
	handlers "github.com/WildEgor/gNotifier/internal/handlers/http"
	middleware "github.com/WildEgor/gNotifier/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type HTTPRouter struct {
	ha                *adapters.HealthCheckAdapter
	storeTokenHandler *handlers.StoreTokenHandler
	unsubTokenHandler *handlers.UnsubTokenHandler
}

func NewHTTPRouter(
	ha *adapters.HealthCheckAdapter,
	storeTokenHandler *handlers.StoreTokenHandler,
	unsubTokenHandler *handlers.UnsubTokenHandler,
) *HTTPRouter {
	return &HTTPRouter{
		ha:                ha,
		storeTokenHandler: storeTokenHandler,
		unsubTokenHandler: unsubTokenHandler,
	}
}

func (r *HTTPRouter) SetupRoutes(app *fiber.App) error {
	hCfg := middleware.HealthCheckConfig{
		Endpoint:           "/api/v1/health/check",
		HealthCheckAdapter: r.ha,
	}

	app.Use(middleware.HealthCheck(&hCfg))

	v1 := app.Group("/api/v1")
	healthCheckController := v1.Group("/health")
	healthCheckController.Get("/check", middleware.HealthCheck(&hCfg))

	tokensController := v1.Group("/tokens")
	tokensController.Post("/store", r.storeTokenHandler.Handle)
	tokensController.Post("/unsub", r.unsubTokenHandler.Handle)

	return nil
}
