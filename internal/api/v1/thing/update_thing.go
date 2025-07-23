package thing

import (
	"context"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/factory"
	"github.com/dmitriygnatenko/homethings/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings/internal/helpers/request"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
)

// @Router 		/api/v1/things/{thingId} [put]
// @Param       thingId path int true "Thing ID"
// @Param       data body dto.UpdateThingRequest true "Request body"
// @Success     200 {object} dto.ThingResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update thing
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateThingHandler(
	tm TransactionManager,
	thingRepository ThingRepository,
	placeThingRepository PlaceThingRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		id, err := request.ConvertToUint64(fctx, "thingId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		req := dto.UpdateThingRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err = validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		err = tm.ReadCommitted(ctx, func(ctx context.Context) error {
			thing, txErr := thingRepository.Get(ctx, id)
			if txErr != nil {
				return txErr
			}

			placeThing, txErr := placeThingRepository.GetByThingID(ctx, id)
			if txErr != nil {
				return txErr
			}

			if req.Title != thing.Title || req.Description != thing.Description {
				txErr = thingRepository.Update(ctx, mappers.ToUpdateThingRequest(id, req))
				if txErr != nil {
					return txErr
				}
			}

			if placeThing.PlaceID != req.PlaceID {
				txErr = placeThingRepository.UpdatePlace(ctx, mappers.ToUpdatePlaceThingRequest(id, req.PlaceID))
				if txErr != nil {
					return txErr
				}
			}

			return nil
		})

		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := thingRepository.Get(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingResponse(*res))
	}
}
