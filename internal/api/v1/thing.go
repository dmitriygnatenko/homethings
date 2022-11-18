package v1

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
	return func(ctx *fiber.Ctx) error {

		return ctx.JSON(dto.ThingResponse{})
	}
}

// @Router 		/v1/places/{id}/things [get]
// @Param       id path int true "Place ID"
// @Success     200 {array} dto.ThingResponse
// @Summary     Get things by place ID
// @Tags  		Places
// @security 	BasicAuth
// @Produce     json
func GetPlaceThingsHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		return ctx.JSON(dto.ThingResponse{})
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
	return func(ctx *fiber.Ctx) error {
		req := dto.AddThingRequest{}
		if err := ctx.BodyParser(&req); err != nil {
			return err
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return helpers.FormatError(err)
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx.Context(), sql.LevelReadCommitted)
		if err != nil {
			return err
		}

		id, err := sp.GetThingRepository().Add(ctx.Context(), mappers.ConvertToAddThingRequestModel(req), tx)
		if err != nil {
			return err
		}

		err = sp.GetPlaceThingRepository().Add(ctx.Context(), models.AddPlaceThingRequest{
			PlaceID: req.PlaceID,
			ThingID: id,
		}, tx)

		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		res, err := sp.GetThingRepository().Get(ctx.Context(), id)
		if err != nil {
			return err
		}

		return ctx.JSON(mappers.ConvertToThingResponseDTO(*res))
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
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return err
		}

		req := dto.UpdatePlaceRequest{}
		if err = ctx.BodyParser(&req); err != nil {
			return err
		}

		_ = id

		return ctx.JSON(dto.ThingResponse{})
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
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return err
		}

		_ = id

		return ctx.JSON(dto.EmptyResponse{})
	}
}
