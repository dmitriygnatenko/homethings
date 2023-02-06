package place

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/places/tree [get]
// @Success     200 {object} dto.PlaceTreeResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get places tree
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlaceTreeHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		res, err := sp.GetPlaceRepository().GetAll(fctx.Context())
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(mappers.ConvertToPlaceTreeResponseDTO(res))
	}
}
