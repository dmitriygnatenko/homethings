package notification

import (
	"database/sql"
	"errors"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/factory"
	"github.com/dmitriygnatenko/homethings/internal/helpers/request"
)

// @Router 		/api/v1/things/notifications/{thingId} [delete]
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete thing notification
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeleteThingNotificationHandler(
	thingNotificationRepository ThingNotificationRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		id, err := request.ConvertToUint64(fctx, "thingId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
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

		if err = thingNotificationRepository.Delete(ctx, id); err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
