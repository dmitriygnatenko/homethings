package place

import (
	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/factory"
	"github.com/dmitriygnatenko/homethings/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings/internal/helpers/request"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
)

// @Router 		/api/v1/places/{placeId} [put]
// @Param       placeId path int true "Place ID"
// @Param       data body dto.UpdatePlaceRequest true "Request body"
// @Success     200 {object} dto.PlaceResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update place
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdatePlaceHandler(placeRepository PlaceRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		id, err := request.ConvertToUint64(fctx, "placeId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		req := dto.UpdatePlaceRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err = validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		err = placeRepository.Update(ctx, mappers.ToUpdatePlaceRequest(id, req))
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := placeRepository.Get(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToPlaceResponse(*res))
	}
}
