package thing

import (
	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things [post]
// @Param       data body dto.AddThingRequest true "Request body"
// @Success     200 {object} dto.ThingResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add thing
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddThingHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddThingRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		id, err := sp.GetThingRepository().Add(ctx, mappers.ConvertToAddThingRequestModel(req), tx)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = sp.GetPlaceThingRepository().Add(ctx, mappers.ConvertToAddPlaceThingRequestModel(id, req.PlaceID), tx)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetThingRepository().CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ConvertToThingResponseDTO(*res))
	}
}
