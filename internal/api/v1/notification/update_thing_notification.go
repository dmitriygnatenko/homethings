package notification

import (
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/location"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/request"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/things/notifications/{thingId} [put]
// @Param       thingId path int true "Thing ID"
// @Param       data body dto.UpdateThingNotificationRequest true "Request body"
// @Success     200 {object} dto.ThingNotificationResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update thing notification
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateThingNotificationHandler(
	thingNotificationRepository ThingNotificationRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := request.ConvertToUint64(fctx, "thingId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		req := dto.UpdateThingNotificationRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err = validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		_, err = thingNotificationRepository.Get(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		dbReq, err := mappers.ToUpdateThingNotificationRequest(id, req)
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err = thingNotificationRepository.Update(ctx, *dbReq); err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := thingNotificationRepository.Get(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingNotificationResponse(*res))
	}
}
