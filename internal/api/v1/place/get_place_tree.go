package place

import (
	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/location"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/places/tree [get]
// @Success     200 {object} dto.PlaceTreeResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get places tree
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlaceTreeHandler(placeRepository PlaceRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		res, err := placeRepository.GetAll(ctx)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToPlaceTreeResponse(res))
	}
}
