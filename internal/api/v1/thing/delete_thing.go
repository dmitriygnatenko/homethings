package thing

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
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
func DeleteThingHandler(
	thingRepository ThingRepository,
	thingTagRepository ThingTagRepository,
	placeThingRepository PlaceThingRepository,
	thingImageRepository ThingImageRepository,
	thingNotificationRepository ThingNotificationRepository,
	fileRepository FileRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("thingId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err = thingRepository.Get(ctx, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tx, err := thingRepository.BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = placeThingRepository.DeleteThing(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		images, err := thingImageRepository.GetByThingID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		imageURLs := make([]string, 0, len(images))
		for i := range images {
			imageURLs = append(imageURLs, images[i].Image)

			if err = thingImageRepository.Delete(ctx, images[i].ID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if err = thingTagRepository.DeleteByThingID(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = thingNotificationRepository.Delete(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = thingRepository.Delete(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = thingRepository.CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(imageURLs) > 0 {
			for i := range imageURLs {
				if err = fileRepository.Delete(imageURLs[i]); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
