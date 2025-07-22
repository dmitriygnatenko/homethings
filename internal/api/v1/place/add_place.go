package place

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i TransactionManager,PlaceRepository,ThingRepository,PlaceImageRepository,ThingImageRepository,PlaceThingRepository,ThingTagRepository,ThingNotificationRepository,FileRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/factory"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings-v1/internal/mappers"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

type (
	TransactionManager interface {
		ReadCommitted(context.Context, func(ctx context.Context) error) error
	}

	PlaceRepository interface {
		GetAll(ctx context.Context) ([]models.Place, error)
		Get(ctx context.Context, id uint64) (*models.Place, error)
		GetNestedPlaces(ctx context.Context, id uint64) ([]models.Place, error)
		Add(ctx context.Context, req models.AddPlaceRequest) (uint64, error)
		Update(ctx context.Context, req models.UpdatePlaceRequest) error
		Delete(ctx context.Context, id uint64) error
	}

	ThingRepository interface {
		GetByPlaceID(ctx context.Context, id uint64) ([]models.Thing, error)
		Delete(ctx context.Context, id uint64) error
	}

	PlaceImageRepository interface {
		GetByPlaceID(ctx context.Context, id uint64) ([]models.Image, error)
		Delete(ctx context.Context, id uint64) error
	}

	ThingImageRepository interface {
		GetByThingID(ctx context.Context, id uint64) ([]models.Image, error)
		Delete(ctx context.Context, id uint64) error
	}

	PlaceThingRepository interface {
		DeleteThing(ctx context.Context, id uint64) error
	}

	ThingTagRepository interface {
		DeleteByThingID(ctx context.Context, id uint64) error
	}

	ThingNotificationRepository interface {
		Delete(ctx context.Context, id uint64) error
	}

	FileRepository interface {
		Delete(path string) error
	}
)

// @Router 		/api/v1/places [post]
// @Param       data body dto.AddPlaceRequest true "Request body"
// @Success     200 {object} dto.PlaceResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add place
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddPlaceHandler(
	placeRepository PlaceRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddPlaceRequest{}

		if err := fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		id, err := placeRepository.Add(ctx, mappers.ToAddPlaceRequest(req))
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := placeRepository.Get(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToPlaceResponse(*res))
	}
}
