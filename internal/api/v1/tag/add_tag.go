package tag

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i TagRepository,ThingRepository,ThingTagRepository,TransactionManager -o ./mocks/ -s "_minimock.go"

import (
	"context"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/location"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

type (
	TransactionManager interface {
		ReadCommitted(context.Context, func(ctx context.Context) error) error
	}

	TagRepository interface {
		GetAll(ctx context.Context) ([]models.Tag, error)
		Get(ctx context.Context, id uint64) (*models.Tag, error)
		GetByThingID(ctx context.Context, id uint64) ([]models.Tag, error)
		Add(ctx context.Context, req models.AddTagRequest) (uint64, error)
		Update(ctx context.Context, req models.UpdateTagRequest) error
		Delete(ctx context.Context, id uint64) error
	}

	ThingRepository interface {
		Get(ctx context.Context, id uint64) (*models.Thing, error)
	}

	ThingTagRepository interface {
		Add(ctx context.Context, req models.AddThingTagRequest) error
		Delete(ctx context.Context, req models.DeleteThingTagRequest) error
		DeleteByTagID(ctx context.Context, id uint64) error
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
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		id, err := tagRepository.Add(ctx, mappers.ToAddTagRequest(req))
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := tagRepository.Get(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToTagResponse(*res))
	}
}
