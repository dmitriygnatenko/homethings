package place

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/places/{id} [delete]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.EmptyResponse
// @Summary     Delete place TODO
// @Tags  		Places
// @security 	BasicAuth
// @Accept      json
// @Produce     json
func DeletePlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return err
		}

		_ = id

		return ctx.JSON(dto.EmptyResponse{})
	}
}
