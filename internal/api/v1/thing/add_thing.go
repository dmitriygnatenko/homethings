package thing

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i ThingRepository,PlaceThingRepository,ThingTagRepository,ThingImageRepository,ThingNotificationRepository,FileRepository,TransactionManager -o ./mocks/ -s "_minimock.go"

import (
	"context"

	"git.dmitriygnatenko.ru/dima/go-common/logger"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/location"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	TransactionManager interface {
		ReadCommitted(context.Context, func(ctx context.Context) error) error
	}

	ThingRepository interface {
		Get(ctx context.Context, id uint64) (*models.Thing, error)
		Search(ctx context.Context, search string) ([]models.Thing, error)
		GetByPlaceID(ctx context.Context, id uint64) ([]models.Thing, error)
		GetAllByPlaceID(ctx context.Context, id uint64) ([]models.Thing, error)
		Add(ctx context.Context, req models.AddThingRequest) (uint64, error)
		Update(ctx context.Context, req models.UpdateThingRequest) error
		Delete(ctx context.Context, id uint64) error
	}

	PlaceThingRepository interface {
		Add(ctx context.Context, req models.AddPlaceThingRequest) error
		GetByThingID(ctx context.Context, id uint64) (*models.PlaceThing, error)
		UpdatePlace(ctx context.Context, req models.UpdatePlaceThingRequest) error
		DeleteThing(ctx context.Context, id uint64) error
	}

	ThingTagRepository interface {
		GetByPlaceID(ctx context.Context, id uint64) ([]models.ThingTag, error)
		DeleteByThingID(ctx context.Context, id uint64) error
	}

	ThingImageRepository interface {
		GetByThingID(ctx context.Context, id uint64) ([]models.Image, error)
		Delete(ctx context.Context, id uint64) error
	}

	ThingNotificationRepository interface {
		Delete(ctx context.Context, id uint64) error
	}

	FileRepository interface {
		Delete(path string) error
	}
)

// @Router 		/api/v1/things [post]
// @Param       data body dto.AddThingRequest true "Request body"
// @Success     200 {object} dto.ThingResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add thing
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddThingHandler(
	tm TransactionManager,
	thingRepository ThingRepository,
	placeThingRepository PlaceThingRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddThingRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		var id uint64

		err := tm.ReadCommitted(ctx, func(ctx context.Context) error {
			var txErr error
			id, txErr = thingRepository.Add(ctx, mappers.ToAddThingRequest(req))
			if txErr != nil {
				return txErr
			}

			txErr = placeThingRepository.Add(ctx, mappers.ToAddPlaceThingRequest(id, req.PlaceID))
			if txErr != nil {
				return txErr
			}

			return nil
		})

		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := thingRepository.Get(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingResponse(*res))
	}
}
