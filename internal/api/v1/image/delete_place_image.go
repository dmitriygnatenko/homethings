package image

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/images/place/{id} [delete]
// @Param       id path int true "Image ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete image
// @Tags  		Images
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeletePlaceImageHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		image, err := sp.GetPlaceImageRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetPlaceImageRepository().Delete(ctx, id, nil); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetFileRepository().Delete(image.Image); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
