package thing

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
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
func GetPlaceThingsHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("placeId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		things, err := sp.GetThingRepository().GetAllByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tags, err := sp.GetThingTagRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		things = helpers.ApplyLocation(fctx, things)
		tags = helpers.ApplyLocation(fctx, tags)

		return fctx.JSON(mappers.ConvertToThingsExtResponseDTO(things, tags))
	}
}
