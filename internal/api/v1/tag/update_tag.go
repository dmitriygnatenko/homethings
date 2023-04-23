package tag

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/tags/{tagId} [put]
// @Param       tagId path int true "Tag ID"
// @Param       data body dto.UpdateTagRequest true "Request body"
// @Success     200 {object} dto.TagResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateTagHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("tagId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		req := dto.UpdateTagRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err = validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		err = sp.GetTagRepository().Update(ctx, mappers.ConvertToUpdateTagRequestModel(id, req), nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := sp.GetTagRepository().Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ConvertToTagResponseDTO(*res))
	}
}
