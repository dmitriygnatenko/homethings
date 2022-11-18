package v1

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	defaultTxLevel = sql.LevelReadCommitted
)

// @Router 		/v1/things/{id} [get]
// @Param       id path int true "Thing ID"
// @Success     200 {object} dto.ThingResponse
// @Failure     404 {object} dto.EmptyResponse
// @Summary     Get one thing by ID
// @Tags  		Things
// @security 	BasicAuth
// @Produce     json
func GetThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return err
		}

		res, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return helpers.NotFoundJsonResponse(dto.EmptyResponse{})
			}
			return err
		}

		return fctx.JSON(mappers.ConvertToThingResponseDTO(*res))
	}
}

// @Router 		/v1/places/{id}/things [get]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.ThingsResponse
// @Summary     Get things by place ID
// @Tags  		Places
// @security 	BasicAuth
// @Produce     json
func GetPlaceThingsHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return err
		}

		res, err := sp.GetThingRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return err
		}

		return fctx.JSON(mappers.ConvertToThingsResponseDTO(res))
	}
}

// @Router 		/v1/things [post]
// @Param       data body dto.AddThingRequest true "Request body"
// @Success     200 {object} dto.ThingResponse
// @Failure     400 {array} dto.ErrorResponse
// @Summary     Add thing
// @Tags  		Things
// @security 	BasicAuth
// @Produce     json
func AddThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddThingRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return err
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return helpers.BadRequestJsonResponse(helpers.FormatValidateErrors(err))
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, defaultTxLevel)
		if err != nil {
			return err
		}

		id, err := sp.GetThingRepository().Add(ctx, mappers.ConvertToAddThingRequestModel(req), tx)
		if err != nil {
			return err
		}

		err = sp.GetPlaceThingRepository().Add(ctx, mappers.ConvertToAddPlaceThingRequestModel(id, req.PlaceID), tx)
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		res, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			return err
		}

		return fctx.JSON(mappers.ConvertToThingResponseDTO(*res))
	}
}

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
			return err
		}

		req := dto.UpdateThingRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			return err
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, defaultTxLevel)
		if err != nil {
			return err
		}

		if req.Title != nil || req.Description != nil {
			err = sp.GetThingRepository().Update(ctx, mappers.ConvertToUpdateThingRequestModel(id, req), tx)
			if err != nil {
				return err
			}
		}

		if req.PlaceID != nil {
			err = sp.GetPlaceThingRepository().UpdatePlace(ctx, mappers.ConvertToUpdatePlaceThingRequestModel(id, *req.PlaceID), tx)
			if err != nil {
				return err
			}
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		res, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			return err
		}

		return fctx.JSON(mappers.ConvertToThingResponseDTO(*res))
	}
}

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
			return err
		}

		_, err = sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return helpers.BadRequestJsonResponse(dto.EmptyResponse{})
			}
			return err
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, defaultTxLevel)
		if err != nil {
			return err
		}

		err = sp.GetPlaceThingRepository().DeleteThing(ctx, id, tx)
		if err != nil {
			return err
		}

		err = sp.GetThingRepository().Delete(ctx, id, tx)
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
