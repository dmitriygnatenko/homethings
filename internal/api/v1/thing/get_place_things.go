package thing

import (
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/things/place/{placeId} [get]
// @Param       placeId path int true "Place ID"
// @Success     200 {object} dto.ThingsExtResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get things by place ID
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlaceThingsHandler(
	thingRepository ThingRepository,
	thingTagRepository ThingTagRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("placeId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		things, err := thingRepository.GetAllByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tags, err := thingTagRepository.GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		things = helpers.ApplyLocation(fctx, things)
		tags = helpers.ApplyLocation(fctx, tags)

		return fctx.JSON(mappers.ToThingsExtResponse(things, tags))
	}
}
