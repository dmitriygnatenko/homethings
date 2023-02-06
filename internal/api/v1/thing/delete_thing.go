package thing

import (
	"database/sql"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/things/{id} [delete]
// @Param       id path int true "Thing ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete thing
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeleteThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		if _, err = sp.GetThingRepository().Get(ctx, id); err != nil {
			if err == sql.ErrNoRows {
				return factory.CreateBadRequestResponse(fctx, nil)
			}
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if err = sp.GetPlaceThingRepository().DeleteThing(ctx, id, tx); err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		images, err := sp.GetThingImageRepository().GetByThingID(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		imageURLs := make([]string, 0, len(images))
		for i := range images {
			imageURLs = append(imageURLs, images[i].Image)

			if err = sp.GetThingImageRepository().Delete(ctx, images[i].ID, tx); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		if err = sp.GetThingRepository().Delete(ctx, id, tx); err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if err = sp.GetThingRepository().CommitTx(tx); err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if len(imageURLs) > 0 {
			for i := range imageURLs {
				if err = sp.GetFileRepository().Delete(imageURLs[i]); err != nil {
					return factory.CreateInternalErrorResponse(fctx, err)
				}
			}
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
