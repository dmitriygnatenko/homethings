package thing

import (
	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things/{id} [put]
// @Param       id path int true "Thing ID"
// @Param       data body dto.UpdateThingRequest true "Request body"
// @Success     200 {object} dto.ThingResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update thing
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateThingHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		req := dto.UpdateThingRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err = validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		thing, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		placeThing, err := sp.GetPlaceThingRepository().GetByThingID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if req.Title != thing.Title || req.Description != thing.Description {
			err = sp.GetThingRepository().Update(ctx, mappers.ConvertToUpdateThingRequestModel(id, req), tx)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if placeThing.PlaceID != req.PlaceID {
			err = sp.GetPlaceThingRepository().UpdatePlace(ctx, mappers.ConvertToUpdatePlaceThingRequestModel(id, req.PlaceID), tx)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if err = sp.GetThingRepository().CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ConvertToThingResponseDTO(*res))
	}
}
