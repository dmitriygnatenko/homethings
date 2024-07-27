package place

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
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
		id, err := fctx.ParamsInt("placeId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		_, err = placeRepository.Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		nestedRes, err := placeRepository.GetNestedPlaces(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(nestedRes) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		placeImages, err := placeImageRepository.GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		things, err := thingRepository.GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		var thingImages []models.Image
		thingIDs := make([]int, 0, len(things))

		for _, thing := range things {
			thingIDs = append(thingIDs, thing.ID)

			thingImagesRes, thingImagesErr := thingImageRepository.GetByThingID(ctx, thing.ID)
			if thingImagesErr != nil {
				return fiber.NewError(fiber.StatusInternalServerError, thingImagesErr.Error())
			}

			thingImages = append(thingImages, thingImagesRes...)
		}

		tx, err := placeRepository.BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		placeImageURLs := make([]string, 0, len(placeImages))
		for i := range placeImages {
			placeImageURLs = append(placeImageURLs, placeImages[i].Image)

			if err = placeImageRepository.Delete(ctx, placeImages[i].ID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		for i := range thingImages {
			if err = thingImageRepository.Delete(ctx, thingImages[i].ID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		for _, thingID := range thingIDs {
			if err = placeThingRepository.DeleteThing(ctx, thingID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			if err = thingTagRepository.DeleteByThingID(ctx, thingID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			if err = thingNotificationRepository.Delete(ctx, thingID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			if err = thingRepository.Delete(ctx, thingID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if err = placeRepository.Delete(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = placeRepository.CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(placeImageURLs) > 0 {
			for i := range placeImageURLs {
				if err = fileRepository.Delete(placeImageURLs[i]); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}
		}

		for i := range thingImages {
			if err = fileRepository.Delete(thingImages[i].Image); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
