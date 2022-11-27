package thing

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/things/{id} [delete]
// @Param       id path int true "Thing ID"
// @Success     200 {object} dto.EmptyResponse
// @Summary     Delete thing
// @Tags  		Things
// @security 	BasicAuth
// @Produce     json
func DeleteThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		_, err = sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return factory.CreateBadRequestResponse(fctx, nil)
			}
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, v1.DefaultTxLevel)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		err = sp.GetPlaceThingRepository().DeleteThing(ctx, id, tx)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		err = sp.GetThingRepository().Delete(ctx, id, tx)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if err = sp.GetThingRepository().CommitTx(tx); err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
