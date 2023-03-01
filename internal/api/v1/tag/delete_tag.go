package tag

import (
	"database/sql"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
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
func DeleteTagHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("tagId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err = sp.GetTagRepository().Get(ctx, id); err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tx, err := sp.GetTagRepository().BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetThingTagRepository().DeleteByThingID(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetTagRepository().Delete(ctx, id, tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = sp.GetTagRepository().CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
