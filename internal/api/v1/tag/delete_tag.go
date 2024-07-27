package tag

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
)

// @Router 		/api/v1/tags/{tagId} [delete]
// @Param       tagId path int true "Tag ID"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Delete tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func DeleteTagHandler(
	tagRepository TagRepository,
	thingTagRepository ThingTagRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("tagId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err = tagRepository.Get(ctx, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tx, err := tagRepository.BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = thingTagRepository.DeleteByTagID(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = tagRepository.Delete(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = tagRepository.CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
