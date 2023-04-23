package notification

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things/notification/{thingId} [get]
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.ThingNotificationResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get thing notification
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetThingNotificationHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("thingId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := sp.GetThingNotificationRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ConvertToThingNotificationResponseDTO(*res))
	}
}
