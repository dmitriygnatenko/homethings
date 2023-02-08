package fiber

import (
	"log"
	"strings"
	"time"

	_ "git.dmitriygnatenko.ru/dima/homethings/docs" //nolint
	authAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/auth"
	imageAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/image"
	placeAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/place"
	thingAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/thing"
	userAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/user"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwt "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
)

const (
	appName            = "homethings"
	staticPath         = "../../web"
	swaggerURI         = "/docs"
	limiterMaxRequests = 30
	limiterExpiration  = 30 * time.Second
	loggerTimeFormat   = "02-01-2006 15:04:05"
)

func Init(sp interfaces.IServiceProvider) (*fiber.App, error) {
	fiberApp := fiber.New(getFiberConfig(sp))

	// Configure web root
	fiberApp.Static("/", staticPath)

	// Configure CORS middleware
	fiberApp.Use(cors.New(getCORSConfig(sp)))

	// Configure recover middleware
	fiberApp.Use(recover.New())

	// Configure limiter middleware
	fiberApp.Use(limiter.New(getLimiterConfig()))

	// Configure JWT middleware
	fiberApp.Use(jwt.New(getJWTConfig(sp)))

	// Swagger
	fiberApp.Get(swaggerURI+"/*", swagger.HandlerDefault)

	// API
	api := fiberApp.Group("/api")
	registerHandlers(api, sp)

	return fiberApp, nil
}

func getFiberConfig(sp interfaces.IServiceProvider) fiber.Config {
	return fiber.Config{
		AppName:               appName,
		DisableStartupMessage: true,
		ErrorHandler:          initErrorHandler(sp),
	}
}

func initErrorHandler(sp interfaces.IServiceProvider) fiber.ErrorHandler {
	return func(fctx *fiber.Ctx, err error) error {
		errCode := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			errCode = e.Code
		}

		if err.Error() == "" {
			if errCode == fiber.StatusInternalServerError {
				log.Println(err)
				// mail
			}

			return fctx.Status(errCode).JSON(factory.CreateEmptyResponse())
		}

		return fctx.Status(errCode).JSON(factory.CreateErrorResponse(err))
	}
}

// nolint
func getJWTConfig(sp interfaces.IServiceProvider) jwt.Config {
	return jwt.Config{
		SigningKey: []byte(sp.GetEnvService().GetJWTSecretKey()),
		ErrorHandler: func(fctx *fiber.Ctx, err error) error {
			return factory.CreateBadRequestResponse(fctx, err)
		},
		Filter: func(fctx *fiber.Ctx) bool {
			method := fctx.Method()
			path := fctx.Path()

			if method != fiber.MethodGet && method != fiber.MethodPost &&
				method != fiber.MethodPut && method != fiber.MethodDelete {
				return true
			}

			if strings.HasPrefix(path, "/api/") && path != "/api/v1/auth/login" {
				return false
			}

			return true
		},
	}
}

func getLimiterConfig() limiter.Config {
	return limiter.Config{
		Max:        limiterMaxRequests,
		Expiration: limiterExpiration,
	}
}

func getCORSConfig(sp interfaces.IServiceProvider) cors.Config {
	return cors.Config{
		AllowOrigins: sp.GetEnvService().GetCORSAllowOrigins(),
		AllowMethods: sp.GetEnvService().GetCORSAllowMethods(),
	}
}

func registerHandlers(r fiber.Router, sp interfaces.IServiceProvider) {
	// Public routes
	r.Post("/v1/auth/login", authAPI.LoginHandler(sp))

	// Protected routes
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

	r.Post("/v1/users", userAPI.AddUserHandler(sp))
	r.Put("/v1/users", userAPI.UpdateUserHandler(sp))
}
