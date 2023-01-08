package fiber

import (
	_ "git.dmitriygnatenko.ru/dima/homethings/docs" //nolint
	authAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/auth"
	imageAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/image"
	placeAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/place"
	thingAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/thing"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
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
			return factory.CreateForbiddenResponse(c, nil)
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
	r.Get("/v1/places", placeAPI.GetPlacesHandler(sp))
	r.Get("/v1/places/tree", placeAPI.GetPlaceTreeHandler(sp))
	r.Get("/v1/places/:id<int>", placeAPI.GetPlaceHandler(sp))
	r.Get("/v1/places/:id<int>/things", placeAPI.GetPlaceThingsHandler(sp))
	r.Post("/v1/places", placeAPI.AddPlaceHandler(sp))
	r.Put("/v1/places/:id<int>", placeAPI.UpdatePlaceHandler(sp))
	r.Delete("/v1/places/:id<int>", placeAPI.DeletePlaceHandler(sp))

	r.Get("/v1/things/:id<int>", thingAPI.GetThingHandler(sp))
	r.Post("/v1/things", thingAPI.AddThingHandler(sp))
	r.Put("/v1/things/:id<int>", thingAPI.UpdateThingHandler(sp))
	r.Delete("/v1/things/:id<int>", thingAPI.DeleteThingHandler(sp))

	r.Post("/v1/images", imageAPI.AddImageHandler(sp))

	r.Get("/v1/auth/check", authAPI.CheckAuthHandler(sp))
}
