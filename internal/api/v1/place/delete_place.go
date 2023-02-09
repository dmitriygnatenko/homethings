package place

import (
	"database/sql"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/places/{id} [delete]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete place
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeletePlaceHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		_, err = sp.GetPlaceRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		nestedRes, err := sp.GetPlaceRepository().GetNestedPlaces(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(nestedRes) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		placeImages, err := sp.GetPlaceImageRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		things, err := sp.GetThingRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		var thingImages []models.Image
		thingIDs := make([]int, 0, len(things))

		for _, thing := range things {
			thingIDs = append(thingIDs, thing.ID)

			thingImagesRes, thingImagesErr := sp.GetThingImageRepository().GetByThingID(ctx, thing.ID)
			if thingImagesErr != nil {
				return fiber.NewError(fiber.StatusInternalServerError, thingImagesErr.Error())
			}

			for i := range thingImagesRes {
				thingImages = append(thingImages, thingImagesRes[i])
			}
		}

		tx, err := sp.GetPlaceRepository().BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		placeImageURLs := make([]string, 0, len(placeImages))
		for i := range placeImages {
			placeImageURLs = append(placeImageURLs, placeImages[i].Image)

			if err = sp.GetPlaceImageRepository().Delete(ctx, placeImages[i].ID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		for i := range thingImages {
			if err = sp.GetThingImageRepository().Delete(ctx, thingImages[i].ID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		for _, thingID := range thingIDs {
			if err = sp.GetPlaceThingRepository().DeleteThing(ctx, thingID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			if err = sp.GetThingRepository().Delete(ctx, thingID, tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if err = sp.GetPlaceRepository().Delete(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetPlaceRepository().CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if len(placeImageURLs) > 0 {
			for i := range placeImageURLs {
				if err = sp.GetFileRepository().Delete(placeImageURLs[i]); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}
		}

		for i := range thingImages {
			if err = sp.GetFileRepository().Delete(thingImages[i].Image); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
