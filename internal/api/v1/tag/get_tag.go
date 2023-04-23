package tag

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/tags/{tagId} [get]
// @Param       tagId path int true "Tag ID"
// @Success     200 {object} dto.TagResponse
// @Failure     404 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get one tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetTagHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("tagId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := sp.GetTagRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusNotFound, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ConvertToTagResponseDTO(*res))
	}
}
