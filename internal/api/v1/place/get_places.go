package place

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/places [get]
// @Success     200 {object} dto.PlacesResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get places list
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlacesHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		res, err := sp.GetPlaceRepository().GetAll(fctx.Context())
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ConvertToPlacesResponseDTO(res))
	}
}
