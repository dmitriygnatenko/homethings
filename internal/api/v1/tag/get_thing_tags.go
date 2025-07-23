package tag

import (
	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings/internal/helpers/request"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
)

// @Router 		/api/v1/tags/thing/{thingId} [get]
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.TagsResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get thing tags
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetThingTagsHandler(tagRepository TagRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		id, err := request.ConvertToUint64(fctx, "thingId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := tagRepository.GetByThingID(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToTagsResponse(res))
	}
}
