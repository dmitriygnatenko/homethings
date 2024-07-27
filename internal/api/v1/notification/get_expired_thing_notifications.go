package notification

import (
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/things/notifications/expired [get]
// @Success     200 {object} dto.ThingNotificationsExtResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get expired thing notifications
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetExpiredThingNotificationsHandler(
	thingNotificationRepository ThingNotificationRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		res, err := thingNotificationRepository.GetExpired(fctx.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingNotificationsExtResponse(res))
	}
}
