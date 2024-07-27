package tag

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i TagRepository,ThingRepository,ThingTagRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

type (
	TagRepository interface {
		GetAll(ctx context.Context) ([]models.Tag, error)
		Get(ctx context.Context, tagID int) (*models.Tag, error)
		GetByThingID(ctx context.Context, thingID int) ([]models.Tag, error)
		Add(ctx context.Context, req models.AddTagRequest, tx *sql.Tx) (int, error)
		Update(ctx context.Context, req models.UpdateTagRequest, tx *sql.Tx) error
		Delete(ctx context.Context, tagID int, tx *sql.Tx) error
		BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
		CommitTx(tx *sql.Tx) error
	}

	ThingRepository interface {
		Get(ctx context.Context, thingID int) (*models.Thing, error)
	}

	ThingTagRepository interface {
		Add(ctx context.Context, req models.AddThingTagRequest, tx *sql.Tx) error
		Delete(ctx context.Context, req models.DeleteThingTagRequest, tx *sql.Tx) error
		DeleteByTagID(ctx context.Context, tagID int, tx *sql.Tx) error
	}
)

// @Router 		/api/v1/tags [post]
// @Param       data body dto.AddTagRequest true "Request body"
// @Success     200 {object} dto.TagResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddTagHandler(tagRepository TagRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddTagRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		id, err := tagRepository.Add(ctx, mappers.ToAddTagRequest(req), nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := tagRepository.Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToTagResponse(*res))
	}
}
