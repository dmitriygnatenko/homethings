package place

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/places/{id} [get]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.PlaceResponse
// @Failure     404 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get one place by ID
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlaceHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := sp.GetPlaceRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ConvertToPlaceResponseDTO(*res))
	}
}
