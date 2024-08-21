package place

import (
	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/location"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/request"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/places/{parentPlaceId}/nested [get]
// @Param       parentPlaceId path int true "Place ID"
// @Success     200 {object} dto.PlacesResponse
// @Failure     500 {object} dto.ErrorResponse
// @Failure     400 {object} dto.ErrorResponse
// @Summary     Get nested places by parent ID
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetNestedPlacesHandler(placeRepository PlaceRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := request.ConvertToUint64(fctx, "parentPlaceId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := placeRepository.GetNestedPlaces(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToPlacesResponse(res))
	}
}
