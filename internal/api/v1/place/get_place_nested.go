package place

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/places/{id}/nested [get]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.PlacesResponse
// @Failure     500 {object} dto.ErrorResponse
// @Failure     400 {object} dto.ErrorResponse
// @Summary     Get nested places by parent ID
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetNestedPlacesHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		res, err := sp.GetPlaceRepository().GetNestedPlaces(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(mappers.ConvertToPlacesResponseDTO(res))
	}
}
