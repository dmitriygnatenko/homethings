package thing

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/things/{id} [put]
// @Param       id path int true "Thing ID"
// @Param       data body dto.UpdateThingRequest true "Request body"
// @Success     200 {object} dto.ThingResponse
// @Summary     Update thing
// @Tags  		Things
// @security 	BasicAuth
// @Produce     json
func UpdateThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		req := dto.UpdateThingRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, v1.DefaultTxLevel)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if req.Title != nil || req.Description != nil {
			err = sp.GetThingRepository().Update(ctx, mappers.ConvertToUpdateThingRequestModel(id, req), tx)
			if err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		if req.PlaceID != nil {
			err = sp.GetPlaceThingRepository().UpdatePlace(ctx, mappers.ConvertToUpdatePlaceThingRequestModel(id, *req.PlaceID), tx)
			if err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		if err = sp.GetThingRepository().CommitTx(tx); err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		res, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(mappers.ConvertToThingResponseDTO(*res))
	}
}
