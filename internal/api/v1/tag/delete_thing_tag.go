package tag

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/tags/{tagId}/thing/{thingId} [delete]
// @Param       tagId path int true "Tag ID"
// @Param       thingId path int true "Thing ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete thing tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeleteThingTagHandler(
	tagRepository TagRepository,
	thingRepository ThingRepository,
	thingTagRepository ThingTagRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		tagID, err := fctx.ParamsInt("tagId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		thingID, err := fctx.ParamsInt("thingId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err = tagRepository.Get(ctx, tagID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if _, err = thingRepository.Get(ctx, thingID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = thingTagRepository.Delete(ctx, mappers.ToDeleteThingTagRequest(tagID, thingID), nil); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
