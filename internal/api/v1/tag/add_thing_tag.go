package tag

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/tags/{id}/thing/{thing_id} [post]
// @Param       id path int true "Tag ID"
// @Param       thing_id path int true "Thing ID"
// @Success     200 {object} dto.TagResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add thing tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddThingTagHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		tagID, err := fctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		thingID, err := fctx.ParamsInt("thing_id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		tag, err := sp.GetTagRepository().Get(ctx, tagID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if _, err = sp.GetThingRepository().Get(ctx, thingID); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetThingTagRepository().Add(ctx, mappers.ConvertToAddThingTagRequestModel(tagID, thingID), nil); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ConvertToTagResponseDTO(*tag))
	}
}
