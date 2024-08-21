package image

import (
	"context"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/request"
)

// @Router 		/api/v1/images/thing/{imageId} [delete]
// @Param       imageId path int true "Image ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete image
// @Tags  		Images
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeleteThingImageHandler(
	tm TransactionManager,
	fileRepository FileRepository,
	thingImageRepository ThingImageRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		id, err := request.ConvertToUint64(fctx, "imageId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		err = tm.ReadCommitted(ctx, func(ctx context.Context) error {
			image, txErr := thingImageRepository.Get(ctx, id)
			if txErr != nil {
				return txErr
			}

			if txErr = thingImageRepository.Delete(ctx, id); txErr != nil {
				return txErr
			}

			if txErr = fileRepository.Delete(image.Image); txErr != nil {
				return txErr
			}

			return nil
		})

		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
