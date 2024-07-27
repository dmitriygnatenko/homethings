package notification

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i ThingNotificationRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
)

type (
	ThingNotificationRepository interface {
		Add(ctx context.Context, req models.AddThingNotificationRequest, tx *sql.Tx) error
		Update(ctx context.Context, req models.UpdateThingNotificationRequest, tx *sql.Tx) error
		Delete(ctx context.Context, thingID int, tx *sql.Tx) error
		Get(ctx context.Context, thingID int) (*models.ThingNotification, error)
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
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		dbReq, err := mappers.ToAddThingNotificationRequest(req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err = thingNotificationRepository.Add(ctx, *dbReq, nil); err != nil {
			if repositories.IsFKViolationError(err) || repositories.IsDuplicateKeyError(err) {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := thingNotificationRepository.Get(ctx, req.ThingID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingNotificationResponse(*res))
	}
}
