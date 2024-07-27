package place

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i PlaceRepository,ThingRepository,PlaceImageRepository,ThingImageRepository,PlaceThingRepository,ThingTagRepository,ThingNotificationRepository,FileRepository -o ./mocks/ -s "_minimock.go"

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
	PlaceRepository interface {
		GetAll(ctx context.Context) ([]models.Place, error)
		Get(ctx context.Context, placeID int) (*models.Place, error)
		GetNestedPlaces(ctx context.Context, placeID int) ([]models.Place, error)
		Add(ctx context.Context, req models.AddPlaceRequest, tx *sql.Tx) (int, error)
		Update(ctx context.Context, req models.UpdatePlaceRequest, tx *sql.Tx) error
		Delete(ctx context.Context, placeID int, tx *sql.Tx) error
		BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
		CommitTx(tx *sql.Tx) error
	}

	ThingRepository interface {
		GetByPlaceID(ctx context.Context, placeID int) ([]models.Thing, error)
		Delete(ctx context.Context, thingID int, tx *sql.Tx) error
	}

	PlaceImageRepository interface {
		GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error)
		Delete(ctx context.Context, imageID int, tx *sql.Tx) error
	}

	ThingImageRepository interface {
		GetByThingID(ctx context.Context, thingID int) ([]models.Image, error)
		Delete(ctx context.Context, imageID int, tx *sql.Tx) error
	}

	PlaceThingRepository interface {
		DeleteThing(ctx context.Context, thingID int, tx *sql.Tx) error
	}

	ThingTagRepository interface {
		DeleteByThingID(ctx context.Context, thingID int, tx *sql.Tx) error
	}

	ThingNotificationRepository interface {
		Delete(ctx context.Context, thingID int, tx *sql.Tx) error
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
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		id, err := placeRepository.Add(ctx, mappers.ToAddPlaceRequest(req), nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := placeRepository.Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToPlaceResponse(*res))
	}
}
