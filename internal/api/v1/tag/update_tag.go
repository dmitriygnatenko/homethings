package tag

import (
	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/factory"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/request"
	"github.com/dmitriygnatenko/homethings-v1/internal/mappers"
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
func UpdateTagHandler(tagRepository TagRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		id, err := request.ConvertToUint64(fctx, "tagId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		req := dto.UpdateTagRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err = validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		err = tagRepository.Update(ctx, mappers.ToUpdateTagRequest(id, req))
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := tagRepository.Get(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToTagResponse(*res))
	}
}
