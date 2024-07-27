package image

import (
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/images/thing/{thingId} [get]
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.ImagesResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get images by thing ID
// @Tags  		Images
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetThingImagesHandler(
	thingImageRepository ThingImageRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("thingId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := thingImageRepository.GetByThingID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToImagesResponse(res))
	}
}
