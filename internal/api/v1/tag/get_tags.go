package tag

import (
	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings-v1/internal/mappers"
)

// @Router 		/api/v1/tags [get]
// @Success     200 {object} dto.TagsResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get tags
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetTagsHandler(tagRepository TagRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		res, err := tagRepository.GetAll(ctx)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToTagsResponse(res))
	}
}
