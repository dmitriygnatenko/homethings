package image

import (
	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/request"
	"github.com/dmitriygnatenko/homethings-v1/internal/mappers"
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

		id, err := request.ConvertToUint64(fctx, "thingId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := thingImageRepository.GetByThingID(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ToImagesResponse(location.ApplyLocation(fctx, res)))
	}
}
