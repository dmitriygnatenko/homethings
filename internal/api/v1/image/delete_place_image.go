package image

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
)

// @Router 		/api/v1/images/place/{imageId} [delete]
// @Param       imageId path int true "Image ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete image
// @Tags  		Images
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeletePlaceImageHandler(
	fileRepository FileRepository,
	placeImageRepository PlaceImageRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("imageId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		image, err := placeImageRepository.Get(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = placeImageRepository.Delete(ctx, id, nil); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = fileRepository.Delete(image.Image); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
