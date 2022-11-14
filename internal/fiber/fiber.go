package fiber

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const (
	appName    = "homethings"
	staticPath = "../../web"
)

func Init(sp interfaces.IServiceProvider) (*fiber.App, error) {
	fiberApp := fiber.New(getConfig(sp))

	// Configure web root
	fiberApp.Static("/", staticPath)

	// Configure CORS middleware
	fiberApp.Use(cors.New())

	// Configure recover middleware
	fiberApp.Use(recover.New())

	// Configure handlers
	fiberApp.Get("/", handler.MainPageHandler(sp))
	fiberApp.Get("/tag/:tag", handler.TagHandler(sp))
	fiberApp.Get("/article/:article", handler.ArticleHandler(sp))

	return fiberApp, nil
}

func getConfig(sp interfaces.IServiceProvider) fiber.Config {
	return fiber.Config{
		AppName:               appName,
		DisableStartupMessage: true,
	}
}
