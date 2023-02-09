package image

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/images/thing/{id} [get]
// @Param       id path int true "Thing ID"
// @Success     200 {object} dto.ImagesResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get images by thing ID
// @Tags  		Images
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetThingImagesHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := sp.GetThingImageRepository().GetByThingID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ConvertToImagesResponseDTO(res))
	}
}
