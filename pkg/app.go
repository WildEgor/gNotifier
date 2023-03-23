package pkg

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
)

var AppSet = wire.NewSet(NewApp)

func NewApp() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	// v1 := app.Group("/api/v1")

	// Server endpoint - sanity check that the server is running
	// statusGroup := v1.Group("/health")

	log.Infof("Application is running on %d port...", 8888)
	return app
}
