package thing

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/things [post]
// @Param       data body dto.AddThingRequest true "Request body"
// @Success     200 {object} dto.ThingResponse
// @Failure     400 {array} dto.ValidateErrorResponse
// @Summary     Add thing
// @Tags  		Things
// @security 	BasicAuth
// @Produce     json
func AddThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddThingRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		tx, err := sp.GetThingRepository().BeginTx(ctx, v1.DefaultTxLevel)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		id, err := sp.GetThingRepository().Add(ctx, mappers.ConvertToAddThingRequestModel(req), tx)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		err = sp.GetPlaceThingRepository().Add(ctx, mappers.ConvertToAddPlaceThingRequestModel(id, req.PlaceID), tx)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
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
