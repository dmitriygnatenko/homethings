package notification

import (
	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
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
		ctx := fctx.Context()
		res, err := thingNotificationRepository.GetExpired(ctx)

		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingNotificationsExtResponse(res))
	}
}
