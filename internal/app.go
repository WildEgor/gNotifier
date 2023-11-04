package app

import (
	"fmt"
	"github.com/WildEgor/gNotifier/internal/configs"

	"github.com/WildEgor/gNotifier/internal/adapters"
	handlers_http "github.com/WildEgor/gNotifier/internal/handlers/http"
	"github.com/WildEgor/gNotifier/internal/repository"
	"github.com/WildEgor/gNotifier/internal/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
)

var AppSet = wire.NewSet(
	NewApp,
	adapters.AdaptersSet,
	repository.RepositoriesSet,
	configs.ConfigSet,
	routers.RoutersSet,
)

func NewApp(
	appConfig *configs.AppConfig,
	httpRouter *routers.HTTPRouter,
	amqpRouter *routers.AMQPRouter,
) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers_http.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(recover.New())

	if !appConfig.IsProduction() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	httpRouter.SetupRoutes(app)
	amqpRouter.SetupRoutes()

	defer amqpRouter.Close()

	log.Info(fmt.Sprintf("Application is running on %v port...", appConfig.Port))

	return app
}
