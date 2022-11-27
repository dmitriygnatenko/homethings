package place

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/places/{id} [put]
// @Param       id path int true "Place ID"
// @Param       data body dto.UpdatePlaceRequest true "Request body"
// @Success     200 {object} dto.PlaceResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update place
// @Tags  		Places
// @security 	BasicAuth
// @Accept      json
// @Produce     json
func UpdatePlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		req := dto.UpdatePlaceRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		var validate = validator.New()
		if err = validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		err = sp.GetPlaceRepository().Update(ctx, mappers.ConvertToUpdatePlaceRequestModel(id, req), nil)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		res, err := sp.GetPlaceRepository().Get(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(mappers.ConvertToPlaceResponseDTO(*res))
	}
}
