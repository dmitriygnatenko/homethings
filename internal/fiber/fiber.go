package fiber

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberJwt "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"

	_ "git.dmitriygnatenko.ru/dima/homethings/docs" //nolint
	authAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/auth"
	imageAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/image"
	notificationAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/notification"
	placeAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/place"
	tagAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/tag"
	thingAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/thing"
	userAPI "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/user"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/middleware/timezone"
	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
	authService "git.dmitriygnatenko.ru/dima/homethings/internal/services/auth"
	envService "git.dmitriygnatenko.ru/dima/homethings/internal/services/env"
	mailerService "git.dmitriygnatenko.ru/dima/homethings/internal/services/mailer"
)

const (
	staticPath = "../../web/public"
	swaggerURI = "/docs/*"
	metricsURI = "/metrics"
)

type (
	ServiceProvider interface {
		GetEnvService() *envService.Service
		GetMailerService() *mailerService.Service
		GetAuthService() *authService.Service
		GetUserRepository() *repositories.UserRepository
		GetPlaceRepository() *repositories.PlaceRepository
		GetThingRepository() *repositories.ThingRepository
		GetTagRepository() *repositories.TagRepository
		GetPlaceThingRepository() *repositories.PlaceThingRepository
		GetPlaceImageRepository() *repositories.PlaceImageRepository
		GetThingImageRepository() *repositories.ThingImageRepository
		GetThingTagRepository() *repositories.ThingTagRepository
		GetThingNotificationRepository() *repositories.ThingNotificationRepository
		GetFileRepository() *repositories.FileRepository
	}
)

func Init(sp ServiceProvider) (*fiber.App, error) {
	fiberApp := fiber.New(getFiberConfig(sp))

	// Configure web root
	fiberApp.Static("/", staticPath)

	// Configure CORS middleware
	fiberApp.Use(cors.New(getCORSConfig(sp)))

	// Configure recover middleware
	fiberApp.Use(recover.New())

	// Configure client timezone middleware
	fiberApp.Use(timezone.New())

	// Configure JWT middleware
	jwtAuth := fiberJwt.New(getJWTConfig(sp))

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

func getFiberConfig(sp ServiceProvider) fiber.Config {
	return fiber.Config{
		AppName:               "Homethings",
		DisableStartupMessage: true,
		ErrorHandler:          getErrorHandler(sp),
	}
}

func getErrorHandler(sp ServiceProvider) fiber.ErrorHandler {
	return func(fctx *fiber.Ctx, err error) error {
		errCode := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
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
func getJWTConfig(sp ServiceProvider) fiberJwt.Config {
	return fiberJwt.Config{
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

func getCORSConfig(sp ServiceProvider) cors.Config {
	return cors.Config{
		AllowOrigins: sp.GetEnvService().GetCORSAllowOrigins(),
		AllowMethods: sp.GetEnvService().GetCORSAllowMethods(),
	}
}

func registerHandlers(r fiber.Router, sp ServiceProvider) {
	// Public routes
	r.Post(
		"/v1/auth/login",
		authAPI.LoginHandler(
			sp.GetAuthService(),
			sp.GetUserRepository(),
		),
	)

	// Protected routes
	r.Get(
		"/v1/auth/check",
		authAPI.CheckAuthHandler(
			sp.GetAuthService(),
			sp.GetUserRepository(),
		),
	)

	r.Get(
		"/v1/places",
		placeAPI.GetPlacesHandler(sp.GetPlaceRepository()),
	)
	r.Get(
		"/v1/places/tree",
		placeAPI.GetPlaceTreeHandler(sp.GetPlaceRepository()),
	)
	r.Get(
		"/v1/places/:placeId<int>",
		placeAPI.GetPlaceHandler(sp.GetPlaceRepository()),
	)
	r.Get(
		"/v1/places/:parentPlaceId<int>/nested",
		placeAPI.GetNestedPlacesHandler(sp.GetPlaceRepository()),
	)
	r.Post(
		"/v1/places",
		placeAPI.AddPlaceHandler(sp.GetPlaceRepository()),
	)
	r.Put(
		"/v1/places/:placeId<int>",
		placeAPI.UpdatePlaceHandler(sp.GetPlaceRepository()),
	)
	r.Delete(
		"/v1/places/:placeId<int>",
		placeAPI.DeletePlaceHandler(
			sp.GetPlaceRepository(),
			sp.GetThingRepository(),
			sp.GetPlaceImageRepository(),
			sp.GetThingImageRepository(),
			sp.GetPlaceThingRepository(),
			sp.GetThingTagRepository(),
			sp.GetThingImageRepository(),
			sp.GetFileRepository(),
		),
	)

	r.Get(
		"/v1/things/:thingId<int>",
		thingAPI.GetThingHandler(sp.GetThingRepository()),
	)
	r.Get(
		"/v1/things/search/:search",
		thingAPI.SearchThingHandler(sp.GetThingRepository()),
	)
	r.Get(
		"/v1/things/place/:placeId<int>",
		thingAPI.GetPlaceThingsHandler(
			sp.GetThingRepository(),
			sp.GetThingTagRepository(),
		),
	)
	r.Post(
		"/v1/things",
		thingAPI.AddThingHandler(
			sp.GetThingRepository(),
			sp.GetPlaceThingRepository(),
		),
	)
	r.Put(
		"/v1/things/:thingId<int>",
		thingAPI.UpdateThingHandler(
			sp.GetThingRepository(),
			sp.GetPlaceThingRepository(),
		),
	)
	r.Delete(
		"/v1/things/:thingId<int>",
		thingAPI.DeleteThingHandler(
			sp.GetThingRepository(),
			sp.GetThingTagRepository(),
			sp.GetPlaceThingRepository(),
			sp.GetThingImageRepository(),
			sp.GetThingNotificationRepository(),
			sp.GetFileRepository(),
		),
	)

	r.Get(
		"/v1/images/place/:placeId<int>",
		imageAPI.GetPlaceImagesHandler(
			sp.GetThingImageRepository(),
			sp.GetPlaceImageRepository(),
		),
	)
	r.Get(
		"/v1/images/thing/:thingId<int>",
		imageAPI.GetThingImagesHandler(sp.GetThingImageRepository()))
	r.Post(
		"/v1/images",
		imageAPI.AddImageHandler(
			sp.GetFileRepository(),
			sp.GetThingImageRepository(),
			sp.GetPlaceImageRepository(),
		),
	)
	r.Delete(
		"/v1/images/place/:imageId<int>",
		imageAPI.DeletePlaceImageHandler(
			sp.GetFileRepository(),
			sp.GetPlaceImageRepository(),
		),
	)
	r.Delete(
		"/v1/images/thing/:imageId<int>",
		imageAPI.DeleteThingImageHandler(
			sp.GetFileRepository(),
			sp.GetThingImageRepository(),
		),
	)

	r.Get(
		"/v1/tags",
		tagAPI.GetTagsHandler(sp.GetTagRepository()),
	)
	r.Get(
		"/v1/tags/:tagId<int>",
		tagAPI.GetTagHandler(sp.GetTagRepository()),
	)
	r.Get(
		"/v1/tags/thing/:thingId<int>",
		tagAPI.GetThingTagsHandler(sp.GetTagRepository()),
	)
	r.Post(
		"/v1/tags",
		tagAPI.AddTagHandler(sp.GetTagRepository()),
	)
	r.Post(
		"/v1/tags/:tagId<int>/thing/:thingId<int>",
		tagAPI.AddThingTagHandler(
			sp.GetTagRepository(),
			sp.GetThingRepository(),
			sp.GetThingTagRepository(),
		),
	)
	r.Put(
		"/v1/tags/:tagId<int>",
		tagAPI.UpdateTagHandler(sp.GetTagRepository()),
	)
	r.Delete(
		"/v1/tags/:tagId<int>",
		tagAPI.DeleteTagHandler(
			sp.GetTagRepository(),
			sp.GetThingTagRepository(),
		),
	)
	r.Delete(
		"/v1/tags/:tagId<int>/thing/:thingId<int>",
		tagAPI.DeleteThingTagHandler(
			sp.GetTagRepository(),
			sp.GetThingRepository(),
			sp.GetThingTagRepository(),
		),
	)
	r.Get(
		"/v1/tags/:id<int>",
		tagAPI.GetTagHandler(sp.GetTagRepository()),
	)
	r.Post(
		"/v1/tags",
		tagAPI.AddTagHandler(sp.GetTagRepository()),
	)
	r.Put(
		"/v1/tags/:id<int>",
		tagAPI.UpdateTagHandler(sp.GetTagRepository()),
	)
	r.Delete(
		"/v1/tags/:id<int>",
		tagAPI.DeleteTagHandler(
			sp.GetTagRepository(),
			sp.GetThingTagRepository(),
		),
	)

	r.Get(
		"/v1/things/notifications/:thingId<int>",
		notificationAPI.GetThingNotificationHandler(sp.GetThingNotificationRepository()),
	)
	r.Get(
		"/v1/things/notifications/expired",
		notificationAPI.GetExpiredThingNotificationsHandler(sp.GetThingNotificationRepository()),
	)
	r.Post(
		"/v1/things/notifications",
		notificationAPI.AddThingNotificationHandler(sp.GetThingNotificationRepository()),
	)
	r.Put(
		"/v1/things/notifications/:thingId<int>",
		notificationAPI.UpdateThingNotificationHandler(sp.GetThingNotificationRepository()),
	)
	r.Delete(
		"/v1/things/notifications/:thingId<int>",
		notificationAPI.DeleteThingNotificationHandler(sp.GetThingNotificationRepository()),
	)

	r.Post(
		"/v1/users",
		userAPI.AddUserHandler(
			sp.GetAuthService(),
			sp.GetUserRepository(),
		),
	)
	r.Put(
		"/v1/users",
		userAPI.UpdateUserHandler(
			sp.GetAuthService(),
			sp.GetUserRepository(),
		),
	)
}
