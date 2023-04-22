package notification

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things/notification/{thingId} [put]
// @Param       thingId path int true "Thing ID"
// @Param       data body dto.UpdateThingNotificationRequest true "Request body"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update thing notification
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateThingNotificationHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		return fctx.JSON(dto.EmptyResponse{})
	}
}
