package notification

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things/notification/{thingId} [delete]
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete thing notification
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeleteThingNotificationHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("thingId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		_, err = sp.GetThingNotificationRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetThingNotificationRepository().Delete(ctx, id, nil); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
