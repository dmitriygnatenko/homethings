package fiber

import (
	_ "git.dmitriygnatenko.ru/dima/homethings/docs" //nolint
	authAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/auth"
	imageAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/image"
	placeAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/place"
	tagAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/tag"
	thingAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/thing"
	userAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/user"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwt "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"log"
)

const (
	staticPath = "../../web"
	swaggerURI = "/docs/*"
	metricsURI = "/metrics"
)

func Init(sp interfaces.ServiceProvider) (*fiber.App, error) {
	fiberApp := fiber.New(getFiberConfig(sp))

	// Configure web root
	fiberApp.Static("/", staticPath)

	// Configure CORS middleware
	fiberApp.Use(cors.New(getCORSConfig(sp)))

	// Configure recover middleware
	fiberApp.Use(recover.New())

	// Configure JWT middleware
	jwtAuth := jwt.New(getJWTConfig(sp))

	// Configure Basic auth
	basicAuth := basicauth.New(basicauth.Config{
		Users: map[string]string{
			sp.GetEnvService().GetBasicAuthUser(): sp.GetEnvService().GetBasicAuthPassword(),
		},
	})

	// Swagger
	fiberApp.Get(swaggerURI, swagger.HandlerDefault)

	// Metrics
	fiberApp.Get(metricsURI, basicAuth, monitor.New(getMetricsConfig()))

	// API
	api := fiberApp.Group("/api", jwtAuth)
	registerHandlers(api, sp)

	return fiberApp, nil
}

func getFiberConfig(sp interfaces.ServiceProvider) fiber.Config {
	return fiber.Config{
		AppName:               "Homethings",
		DisableStartupMessage: true,
		ErrorHandler:          getErrorHandler(sp),
	}
}

func getErrorHandler(sp interfaces.ServiceProvider) fiber.ErrorHandler {
	return func(fctx *fiber.Ctx, err error) error {
		errCode := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			errCode = e.Code
		}

		if err.Error() != "" {
			errorsEmail := sp.GetEnvService().GetErrorsEmail()

			if errCode == fiber.StatusInternalServerError && errorsEmail != "" {
				log.Println(err)
				// nolint
				sp.GetMailerService().Send(
					errorsEmail,
					"AUTO - Homethings error",
					err.Error(),
				)
			}

			return fctx.Status(errCode).JSON(factory.CreateErrorResponse(err))
		}

		return fctx.Status(errCode).JSON(factory.CreateEmptyResponse())
	}
}

// nolint
func getJWTConfig(sp interfaces.ServiceProvider) jwt.Config {
	return jwt.Config{
		SigningKey: []byte(sp.GetEnvService().GetJWTSecretKey()),
		ErrorHandler: func(fctx *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		},
		Filter: func(fctx *fiber.Ctx) bool {
			method := fctx.Method()
			path := fctx.Path()

			if method != fiber.MethodGet && method != fiber.MethodPost &&
				method != fiber.MethodPut && method != fiber.MethodDelete {
				return true
			}

			if path == "/api/v1/auth/login" {
				return true
			}

			return false
		},
	}
}

func getMetricsConfig() monitor.Config {
	return monitor.Config{
		Title: "Homethings metrics",
	}
}

func getCORSConfig(sp interfaces.ServiceProvider) cors.Config {
	return cors.Config{
		AllowOrigins: sp.GetEnvService().GetCORSAllowOrigins(),
		AllowMethods: sp.GetEnvService().GetCORSAllowMethods(),
	}
}

func registerHandlers(r fiber.Router, sp interfaces.ServiceProvider) {
	// Public routes
	r.Post("/v1/auth/login", authAPI.LoginHandler(sp))

	// Protected routes
	r.Get("/v1/auth/check", authAPI.CheckAuthHandler(sp))

	r.Get("/v1/places", placeAPI.GetPlacesHandler(sp))
	r.Get("/v1/places/tree", placeAPI.GetPlaceTreeHandler(sp))
	r.Get("/v1/places/:id<int>", placeAPI.GetPlaceHandler(sp))
	r.Get("/v1/places/:id<int>/things", placeAPI.GetPlaceThingsHandler(sp))
	r.Get("/v1/places/:id<int>/nested", placeAPI.GetNestedPlacesHandler(sp))
	r.Post("/v1/places", placeAPI.AddPlaceHandler(sp))
	r.Put("/v1/places/:id<int>", placeAPI.UpdatePlaceHandler(sp))
	r.Delete("/v1/places/:id<int>", placeAPI.DeletePlaceHandler(sp))

	r.Get("/v1/things/:id<int>", thingAPI.GetThingHandler(sp))
	r.Get("/v1/things/search/:search", thingAPI.SearchThingHandler(sp))
	r.Post("/v1/things", thingAPI.AddThingHandler(sp))
	r.Put("/v1/things/:id<int>", thingAPI.UpdateThingHandler(sp))
	r.Delete("/v1/things/:id<int>", thingAPI.DeleteThingHandler(sp))

	r.Get("/v1/images/place/:id<int>", imageAPI.GetPlaceImagesHandler(sp))
	r.Get("/v1/images/thing/:id<int>", imageAPI.GetThingImagesHandler(sp))
	r.Post("/v1/images", imageAPI.AddImageHandler(sp))
	r.Delete("/v1/images/place/:id<int>", imageAPI.DeletePlaceImageHandler(sp))
	r.Delete("/v1/images/thing/:id<int>", imageAPI.DeleteThingImageHandler(sp))

	r.Get("/v1/tags/:id<int>", tagAPI.GetTagHandler(sp))
	r.Post("/v1/tags", tagAPI.AddTagHandler(sp))
	r.Put("/v1/tags/:id<int>", tagAPI.UpdateTagHandler(sp))
	r.Delete("/v1/tags/:id<int>", tagAPI.DeleteTagHandler(sp))

	r.Post("/v1/users", userAPI.AddUserHandler(sp))
	r.Put("/v1/users", userAPI.UpdateUserHandler(sp))
}
