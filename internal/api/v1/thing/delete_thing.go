package thing

import (
	"database/sql"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things/{thingId} [delete]
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete thing
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeleteThingHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("thingId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err = sp.GetThingRepository().Get(ctx, id); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetPlaceThingRepository().DeleteThing(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		images, err := sp.GetThingImageRepository().GetByThingID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		imageURLs := make([]string, 0, len(images))
		for i := range images {
			imageURLs = append(imageURLs, images[i].Image)

			if err = sp.GetThingImageRepository().Delete(ctx, images[i].ID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if err = sp.GetThingTagRepository().DeleteByThingID(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetThingRepository().Delete(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetThingRepository().CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(imageURLs) > 0 {
			for i := range imageURLs {
				if err = sp.GetFileRepository().Delete(imageURLs[i]); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
