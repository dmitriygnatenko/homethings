package v1

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/places/{id}/things [get]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.ThingsResponse
// @Summary     Get things by place ID
// @Tags  		Places
// @security 	BasicAuth
// @Produce     json
func GetPlaceThingsHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		res, err := sp.GetThingRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(mappers.ConvertToThingsResponseDTO(res))
	}
}
