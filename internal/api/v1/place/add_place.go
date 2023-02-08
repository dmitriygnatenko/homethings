package place

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/places [post]
// @Param       data body dto.AddPlaceRequest true "Request body"
// @Success     200 {object} dto.PlaceResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add place
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddPlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddPlaceRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		id, err := sp.GetPlaceRepository().Add(ctx, mappers.ConvertToAddPlaceRequestModel(req), nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := sp.GetPlaceRepository().Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ConvertToPlaceResponseDTO(*res))
	}
}
