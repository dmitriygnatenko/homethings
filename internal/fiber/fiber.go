package fiber

import (
	"errors"

	"git.dmitriygnatenko.ru/dima/go-common/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberJwt "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"

	_ "git.dmitriygnatenko.ru/dima/homethings/docs" // nolint
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
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/auth"
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/config"
)

const (
	staticPath = "../../web/public"
	swaggerURI = "/docs/*"
	metricsURI = "/metrics"
)

type (
	ServiceProvider interface {
		ConfigService() *config.Service
		AuthService() *auth.Service
		TransactionManager() *db.TxManager
		UserRepository() *repositories.UserRepository
		PlaceRepository() *repositories.PlaceRepository
		ThingRepository() *repositories.ThingRepository
		TagRepository() *repositories.TagRepository
		PlaceThingRepository() *repositories.PlaceThingRepository
		PlaceImageRepository() *repositories.PlaceImageRepository
		ThingImageRepository() *repositories.ThingImageRepository
		ThingTagRepository() *repositories.ThingTagRepository
		ThingNotificationRepository() *repositories.ThingNotificationRepository
		FileRepository() *repositories.FileRepository
	}
)

func Init(sp ServiceProvider) (*fiber.App, error) {
	fiberApp := fiber.New(getFiberConfig())

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
			sp.ConfigService().BasicAuthUser(): sp.ConfigService().BasicAuthPassword(),
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

func getFiberConfig() fiber.Config {
	return fiber.Config{
		AppName:               "Homethings",
		DisableStartupMessage: true,
		ErrorHandler:          getErrorHandler(),
	}
}

func getErrorHandler() fiber.ErrorHandler {
	return func(fctx *fiber.Ctx, err error) error {
		errCode := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			errCode = e.Code
		}

		if err.Error() != "" {
			return fctx.Status(errCode).JSON(factory.CreateErrorResponse(err))
		}

		return fctx.Status(errCode).JSON(factory.CreateEmptyResponse())
	}
}

// nolint
func getJWTConfig(sp ServiceProvider) fiberJwt.Config {
	return fiberJwt.Config{
		SigningKey: []byte(sp.ConfigService().JWTSecretKey()),
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
		AllowOrigins: sp.ConfigService().CORSAllowOrigins(),
		AllowMethods: sp.ConfigService().CORSAllowMethods(),
	}
}

func registerHandlers(r fiber.Router, sp ServiceProvider) {
	// Public routes
	r.Post(
		"/v1/auth/login",
		authAPI.LoginHandler(
			sp.AuthService(),
			sp.UserRepository(),
		),
	)

	// Protected routes
	r.Get(
		"/v1/auth/check",
		authAPI.CheckAuthHandler(
			sp.AuthService(),
			sp.UserRepository(),
		),
	)

	r.Get(
		"/v1/places",
		placeAPI.GetPlacesHandler(
			sp.PlaceRepository(),
		),
	)

	r.Get(
		"/v1/places/tree",
		placeAPI.GetPlaceTreeHandler(
			sp.PlaceRepository(),
		),
	)

	r.Get(
		"/v1/places/:placeId<int>",
		placeAPI.GetPlaceHandler(
			sp.PlaceRepository(),
		),
	)

	r.Get(
		"/v1/places/:parentPlaceId<int>/nested",
		placeAPI.GetNestedPlacesHandler(
			sp.PlaceRepository(),
		),
	)

	r.Post(
		"/v1/places",
		placeAPI.AddPlaceHandler(
			sp.PlaceRepository(),
		),
	)

	r.Put(
		"/v1/places/:placeId<int>",
		placeAPI.UpdatePlaceHandler(
			sp.PlaceRepository(),
		),
	)

	r.Delete(
		"/v1/places/:placeId<int>",
		placeAPI.DeletePlaceHandler(
			sp.TransactionManager(),
			sp.PlaceRepository(),
			sp.ThingRepository(),
			sp.PlaceImageRepository(),
			sp.ThingImageRepository(),
			sp.PlaceThingRepository(),
			sp.ThingTagRepository(),
			sp.ThingImageRepository(),
			sp.FileRepository(),
		),
	)

	r.Get(
		"/v1/things/:thingId<int>",
		thingAPI.GetThingHandler(sp.ThingRepository()),
	)

	r.Get(
		"/v1/things/search/:search",
		thingAPI.SearchThingHandler(sp.ThingRepository()),
	)

	r.Get(
		"/v1/things/place/:placeId<int>",
		thingAPI.GetPlaceThingsHandler(
			sp.ThingRepository(),
			sp.ThingTagRepository(),
		),
	)

	r.Post(
		"/v1/things",
		thingAPI.AddThingHandler(
			sp.TransactionManager(),
			sp.ThingRepository(),
			sp.PlaceThingRepository(),
		),
	)

	r.Put(
		"/v1/things/:thingId<int>",
		thingAPI.UpdateThingHandler(
			sp.TransactionManager(),
			sp.ThingRepository(),
			sp.PlaceThingRepository(),
		),
	)

	r.Delete(
		"/v1/things/:thingId<int>",
		thingAPI.DeleteThingHandler(
			sp.TransactionManager(),
			sp.ThingRepository(),
			sp.ThingTagRepository(),
			sp.PlaceThingRepository(),
			sp.ThingImageRepository(),
			sp.ThingNotificationRepository(),
			sp.FileRepository(),
		),
	)

	r.Get(
		"/v1/images/place/:placeId<int>",
		imageAPI.GetPlaceImagesHandler(
			sp.ThingImageRepository(),
			sp.PlaceImageRepository(),
		),
	)

	r.Get(
		"/v1/images/thing/:thingId<int>",
		imageAPI.GetThingImagesHandler(sp.ThingImageRepository()))

	r.Post(
		"/v1/images",
		imageAPI.AddImageHandler(
			sp.TransactionManager(),
			sp.FileRepository(),
			sp.ThingImageRepository(),
			sp.PlaceImageRepository(),
		),
	)

	r.Delete(
		"/v1/images/place/:imageId<int>",
		imageAPI.DeletePlaceImageHandler(
			sp.TransactionManager(),
			sp.FileRepository(),
			sp.PlaceImageRepository(),
		),
	)

	r.Delete(
		"/v1/images/thing/:imageId<int>",
		imageAPI.DeleteThingImageHandler(
			sp.TransactionManager(),
			sp.FileRepository(),
			sp.ThingImageRepository(),
		),
	)

	r.Get(
		"/v1/tags",
		tagAPI.GetTagsHandler(
			sp.TagRepository(),
		),
	)

	r.Get(
		"/v1/tags/:tagId<int>",
		tagAPI.GetTagHandler(
			sp.TagRepository(),
		),
	)

	r.Get(
		"/v1/tags/thing/:thingId<int>",
		tagAPI.GetThingTagsHandler(
			sp.TagRepository(),
		),
	)

	r.Post(
		"/v1/tags",
		tagAPI.AddTagHandler(
			sp.TagRepository(),
		),
	)

	r.Post(
		"/v1/tags/:tagId<int>/thing/:thingId<int>",
		tagAPI.AddThingTagHandler(
			sp.TagRepository(),
			sp.ThingRepository(),
			sp.ThingTagRepository(),
		),
	)

	r.Put(
		"/v1/tags/:tagId<int>",
		tagAPI.UpdateTagHandler(
			sp.TagRepository(),
		),
	)

	r.Delete(
		"/v1/tags/:tagId<int>",
		tagAPI.DeleteTagHandler(
			sp.TransactionManager(),
			sp.TagRepository(),
			sp.ThingTagRepository(),
		),
	)

	r.Delete(
		"/v1/tags/:tagId<int>/thing/:thingId<int>",
		tagAPI.DeleteThingTagHandler(
			sp.TagRepository(),
			sp.ThingRepository(),
			sp.ThingTagRepository(),
		),
	)

	r.Get(
		"/v1/tags/:id<int>",
		tagAPI.UpdateTagHandler(
			sp.TagRepository(),
		),
	)

	r.Post(
		"/v1/tags",
		tagAPI.UpdateTagHandler(
			sp.TagRepository(),
		),
	)

	r.Put(
		"/v1/tags/:id<int>",
		tagAPI.UpdateTagHandler(
			sp.TagRepository(),
		),
	)

	r.Delete(
		"/v1/tags/:id<int>",
		tagAPI.DeleteTagHandler(
			sp.TransactionManager(),
			sp.TagRepository(),
			sp.ThingTagRepository(),
		),
	)

	r.Get(
		"/v1/things/notifications/:thingId<int>",
		notificationAPI.GetThingNotificationHandler(
			sp.ThingNotificationRepository(),
		),
	)

	r.Get(
		"/v1/things/notifications/expired",
		notificationAPI.GetExpiredThingNotificationsHandler(
			sp.ThingNotificationRepository(),
		),
	)

	r.Post(
		"/v1/things/notifications",
		notificationAPI.AddThingNotificationHandler(
			sp.ThingNotificationRepository(),
		),
	)

	r.Put(
		"/v1/things/notifications/:thingId<int>",
		notificationAPI.UpdateThingNotificationHandler(
			sp.ThingNotificationRepository(),
		),
	)

	r.Delete(
		"/v1/things/notifications/:thingId<int>",
		notificationAPI.DeleteThingNotificationHandler(
			sp.ThingNotificationRepository(),
		),
	)

	r.Post(
		"/v1/users",
		userAPI.AddUserHandler(
			sp.AuthService(),
			sp.UserRepository(),
		),
	)

	r.Put(
		"/v1/users",
		userAPI.UpdateUserHandler(
			sp.AuthService(),
			sp.UserRepository(),
		),
	)
}
