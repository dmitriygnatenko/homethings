package fiber

import (
	apiV1 "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	_ "git.dmitriygnatenko.ru/dima/homethings/internal/docs" //nolint
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

const (
	appName    = "homethings"
	staticPath = "../../web"
	swaggerURI = "/docs"
)

func Init(sp interfaces.IServiceProvider) (*fiber.App, error) {
	fiberApp := fiber.New(getFiberConfig())

	// Configure web root
	fiberApp.Static("/", staticPath)

	// Configure CORS middleware
	fiberApp.Use(cors.New(getCORSConfig(sp)))

	// Configure recover middleware
	fiberApp.Use(recover.New())

	// Swagger
	fiberApp.Get(swaggerURI+"/*", swagger.HandlerDefault)

	// API
	api := fiberApp.Group("/api", basicauth.New(getBasicAuthConfig(sp)))
	registerHandlers(api, sp)

	return fiberApp, nil
}

func getFiberConfig() fiber.Config {
	return fiber.Config{
		AppName:               appName,
		DisableStartupMessage: true,
	}
}

func getBasicAuthConfig(sp interfaces.IServiceProvider) basicauth.Config {
	user := sp.GetEnvService().GetAuthUser()
	password := sp.GetEnvService().GetAuthPassword()

	return basicauth.Config{
		Users: map[string]string{
			user: password,
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusForbidden)
		},
	}
}

func getCORSConfig(sp interfaces.IServiceProvider) cors.Config {
	return cors.Config{
		AllowOrigins: sp.GetEnvService().GetCORSAllowOrigins(),
		AllowMethods: sp.GetEnvService().GetCORSAllowMethods(),
	}
}

func registerHandlers(r fiber.Router, sp interfaces.IServiceProvider) {
	r.Get("/v1/tags", apiV1.GetTagsHandler(sp))

	r.Get("/v1/places/:id<int>", apiV1.GetPlaceHandler(sp))
	r.Get("/v1/places/:id<int>/things", apiV1.GetPlaceThingsHandler(sp))
	r.Post("/v1/places", apiV1.AddPlaceHandler(sp))
	r.Put("/v1/places/:id<int>", apiV1.UpdatePlaceHandler(sp))
	r.Delete("/v1/places/:id<int>", apiV1.DeletePlaceHandler(sp))

	r.Get("/v1/things/:id<int>", apiV1.GetThingHandler(sp))
	r.Post("/v1/things", apiV1.AddThingHandler(sp))
	r.Put("/v1/things/:id<int>", apiV1.UpdateThingHandler(sp))
	r.Delete("/v1/things/:id<int>", apiV1.DeleteThingHandler(sp))
}
