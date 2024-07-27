package place

import (
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
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
		res, err := placeRepository.GetAll(fctx.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToPlaceTreeResponse(res))
	}
}
