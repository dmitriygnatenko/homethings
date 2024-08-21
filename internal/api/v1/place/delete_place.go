package place

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/request"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

// @Router 		/api/v1/places/{placeId} [delete]
// @Param       placeId path int true "Place ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete place
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeletePlaceHandler(
	tm TransactionManager,
	placeRepository PlaceRepository,
	thingRepository ThingRepository,
	placeImageRepository PlaceImageRepository,
	thingImageRepository ThingImageRepository,
	placeThingRepository PlaceThingRepository,
	thingTagRepository ThingTagRepository,
	thingNotificationRepository ThingNotificationRepository,
	fileRepository FileRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := request.ConvertToUint64(fctx, "placeId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		_, err = placeRepository.Get(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		nestedRes, err := placeRepository.GetNestedPlaces(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(nestedRes) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		err = tm.ReadCommitted(ctx, func(ctx context.Context) error {
			placeImages, txErr := placeImageRepository.GetByPlaceID(ctx, id)
			if txErr != nil {
				return txErr
			}

			things, txErr := thingRepository.GetByPlaceID(ctx, id)
			if txErr != nil {
				return txErr
			}

			var thingImages []models.Image
			thingIDs := make([]uint64, 0, len(things))

			for _, thing := range things {
				thingIDs = append(thingIDs, thing.ID)

				thingImagesRes, thingImagesErr := thingImageRepository.GetByThingID(ctx, thing.ID)
				if thingImagesErr != nil {
					return thingImagesErr
				}

				thingImages = append(thingImages, thingImagesRes...)
			}

			placeImageURLs := make([]string, 0, len(placeImages))
			for i := range placeImages {
				placeImageURLs = append(placeImageURLs, placeImages[i].Image)
				if txErr = placeImageRepository.Delete(ctx, placeImages[i].ID); txErr != nil {
					return txErr
				}
			}

			for i := range thingImages {
				if txErr = thingImageRepository.Delete(ctx, thingImages[i].ID); txErr != nil {
					return txErr
				}
			}

			for _, thingID := range thingIDs {
				if txErr = placeThingRepository.DeleteThing(ctx, thingID); txErr != nil {
					return txErr
				}

				if txErr = thingTagRepository.DeleteByThingID(ctx, thingID); txErr != nil {
					return txErr
				}

				if txErr = thingNotificationRepository.Delete(ctx, thingID); txErr != nil {
					return txErr
				}

				if txErr = thingRepository.Delete(ctx, thingID); txErr != nil {
					return txErr
				}
			}

			if txErr = placeRepository.Delete(ctx, id); txErr != nil {
				return txErr
			}

			if len(placeImageURLs) > 0 {
				for i := range placeImageURLs {
					if txErr = fileRepository.Delete(placeImageURLs[i]); txErr != nil {
						return txErr
					}
				}
			}

			for i := range thingImages {
				if txErr = fileRepository.Delete(thingImages[i].Image); txErr != nil {
					return txErr
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
