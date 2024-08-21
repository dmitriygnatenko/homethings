package tag

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/request"
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
	tm TransactionManager,
	tagRepository TagRepository,
	thingTagRepository ThingTagRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := request.ConvertToUint64(fctx, "tagId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if _, err = tagRepository.Get(ctx, id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = tm.ReadCommitted(ctx, func(ctx context.Context) error {
			if txErr := thingTagRepository.DeleteByTagID(ctx, id); txErr != nil {
				return txErr
			}

			if txErr := tagRepository.Delete(ctx, id); txErr != nil {
				return txErr
			}

			return nil
		})

		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
