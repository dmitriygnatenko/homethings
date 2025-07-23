package notification

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i ThingNotificationRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/factory"
	"github.com/dmitriygnatenko/homethings/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
	"github.com/dmitriygnatenko/homethings/internal/models"
	"github.com/dmitriygnatenko/homethings/internal/repositories"
)

type (
	ThingNotificationRepository interface {
		Add(ctx context.Context, req models.AddThingNotificationRequest) error
		Update(ctx context.Context, req models.UpdateThingNotificationRequest) error
		Delete(ctx context.Context, id uint64) error
		Get(ctx context.Context, id uint64) (*models.ThingNotification, error)
		GetExpired(ctx context.Context) ([]models.ExtThingNotification, error)
	}
)

// @Router 		/api/v1/things/notifications [post]
// @Param       data body dto.AddThingNotificationRequest true "Request body"
// @Success     200 {object} dto.ThingNotificationResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add thing notification
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddThingNotificationHandler(
	thingNotificationRepository ThingNotificationRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		req := dto.AddThingNotificationRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		dbReq, err := mappers.ToAddThingNotificationRequest(req)
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err = thingNotificationRepository.Add(ctx, *dbReq); err != nil {
			if repositories.IsFKViolationError(err) || repositories.IsDuplicateKeyError(err) {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			logger.Error(ctx, err.Error())

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := thingNotificationRepository.Get(ctx, req.ThingID)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingNotificationResponse(*res))
	}
}
