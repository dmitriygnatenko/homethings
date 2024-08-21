package thing

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/request"
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
	tm TransactionManager,
	thingRepository ThingRepository,
	thingTagRepository ThingTagRepository,
	placeThingRepository PlaceThingRepository,
	thingImageRepository ThingImageRepository,
	thingNotificationRepository ThingNotificationRepository,
	fileRepository FileRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := request.ConvertToUint64(fctx, "thingId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err = thingRepository.Get(ctx, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = tm.ReadCommitted(ctx, func(ctx context.Context) error {
			if txErr := placeThingRepository.DeleteThing(ctx, id); txErr != nil {
				return txErr
			}

			images, txErr := thingImageRepository.GetByThingID(ctx, id)
			if txErr != nil {
				return txErr
			}

			imageURLs := make([]string, 0, len(images))
			for i := range images {
				imageURLs = append(imageURLs, images[i].Image)

				if txErr = thingImageRepository.Delete(ctx, images[i].ID); txErr != nil {
					return txErr
				}
			}

			if txErr = thingTagRepository.DeleteByThingID(ctx, id); txErr != nil {
				return txErr
			}

			if txErr = thingNotificationRepository.Delete(ctx, id); txErr != nil {
				return txErr
			}

			if txErr = thingRepository.Delete(ctx, id); txErr != nil {
				return txErr
			}

			if len(imageURLs) > 0 {
				for i := range imageURLs {
					if txErr = fileRepository.Delete(imageURLs[i]); txErr != nil {
						return txErr
					}
				}
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
