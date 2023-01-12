package place

import (
	"database/sql"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/places/{id} [delete]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete place
// @Tags  		Places
// @security 	BasicAuth
// @Accept      json
// @Produce     json
func DeletePlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		_, err = sp.GetPlaceRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return factory.CreateBadRequestResponse(fctx, nil)
			}

			return factory.CreateInternalErrorResponse(fctx, err)
		}

		nestedRes, err := sp.GetPlaceRepository().GetNestedPlaces(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if len(nestedRes) > 0 {
			return factory.CreateBadRequestResponse(fctx, nil)
		}

		placeImages, err := sp.GetPlaceImageRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		things, err := sp.GetThingRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		var thingImages []models.Image
		thingIDs := make([]int, 0, len(things))

		for _, thing := range things {
			thingIDs = append(thingIDs, thing.ID)

			thingImagesRes, thingImagesErr := sp.GetThingImageRepository().GetByThingID(ctx, thing.ID)
			if thingImagesErr != nil {
				return factory.CreateInternalErrorResponse(fctx, thingImagesErr)
			}

			for i := range thingImagesRes {
				thingImages = append(thingImages, thingImagesRes[i])
			}
		}

		tx, err := sp.GetPlaceRepository().BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		placeImageURLs := make([]string, 0, len(placeImages))
		for i := range placeImages {
			placeImageURLs = append(placeImageURLs, placeImages[i].Image)

			if err = sp.GetPlaceImageRepository().Delete(ctx, placeImages[i].ID, tx); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		for i := range thingImages {
			if err = sp.GetThingImageRepository().Delete(ctx, thingImages[i].ID, tx); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		for _, thingID := range thingIDs {
			if err = sp.GetPlaceThingRepository().DeleteThing(ctx, thingID, tx); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}

			if err = sp.GetThingRepository().Delete(ctx, thingID, tx); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		if err = sp.GetPlaceRepository().Delete(ctx, id, tx); err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if err = sp.GetPlaceRepository().CommitTx(tx); err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if len(placeImageURLs) > 0 {
			for i := range placeImageURLs {
				if err = sp.GetFileRepository().Delete(placeImageURLs[i]); err != nil {
					return factory.CreateInternalErrorResponse(fctx, err)
				}
			}
		}

		for i := range thingImages {
			if err = sp.GetFileRepository().Delete(thingImages[i].Image); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
