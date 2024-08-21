package tag

import (
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/location"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/request"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/tags/{tagId}/thing/{thingId} [post]
// @Param       tagId path int true "Tag ID"
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.TagResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add thing tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddThingTagHandler(
	tagRepository TagRepository,
	thingRepository ThingRepository,
	thingTagRepository ThingTagRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()

		tagID, err := request.ConvertToUint64(fctx, "tagId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		thingID, err := request.ConvertToUint64(fctx, "thingId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		tag, err := tagRepository.Get(ctx, tagID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if _, err = thingRepository.Get(ctx, thingID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = thingTagRepository.Add(ctx, mappers.ToAddThingTagRequest(tagID, thingID)); err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tag = location.ApplyLocation(fctx, tag)

		return fctx.JSON(mappers.ToTagResponse(*tag))
	}
}
