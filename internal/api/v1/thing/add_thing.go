package thing

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i ThingRepository,PlaceThingRepository,ThingTagRepository,ThingImageRepository,ThingNotificationRepository,FileRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ThingRepository interface {
		Get(ctx context.Context, thingID int) (*models.Thing, error)
		Search(ctx context.Context, search string) ([]models.Thing, error)
		GetByPlaceID(ctx context.Context, placeID int) ([]models.Thing, error)
		GetAllByPlaceID(ctx context.Context, placeID int) ([]models.Thing, error)
		Add(ctx context.Context, req models.AddThingRequest, tx *sql.Tx) (int, error)
		Update(ctx context.Context, req models.UpdateThingRequest, tx *sql.Tx) error
		Delete(ctx context.Context, thingID int, tx *sql.Tx) error
		BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
		CommitTx(tx *sql.Tx) error
	}

	PlaceThingRepository interface {
		Add(ctx context.Context, req models.AddPlaceThingRequest, tx *sql.Tx) error
		GetByThingID(ctx context.Context, thingID int) (*models.PlaceThing, error)
		UpdatePlace(ctx context.Context, req models.UpdatePlaceThingRequest, tx *sql.Tx) error
		DeleteThing(ctx context.Context, thingID int, tx *sql.Tx) error
	}

	ThingTagRepository interface {
		GetByPlaceID(ctx context.Context, placeID int) ([]models.ThingTag, error)
		DeleteByThingID(ctx context.Context, thingID int, tx *sql.Tx) error
	}

	ThingImageRepository interface {
		GetByThingID(ctx context.Context, thingID int) ([]models.Image, error)
		Delete(ctx context.Context, imageID int, tx *sql.Tx) error
	}

	ThingNotificationRepository interface {
		Delete(ctx context.Context, thingID int, tx *sql.Tx) error
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
	thingRepository ThingRepository,
	placeThingRepository PlaceThingRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddThingRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		tx, err := thingRepository.BeginTx(ctx, API.DefaultTxLevel)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		id, err := thingRepository.Add(ctx, mappers.ToAddThingRequest(req), tx)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = placeThingRepository.Add(ctx, mappers.ToAddPlaceThingRequest(id, req.PlaceID), tx)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if err = thingRepository.CommitTx(tx); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res, err := thingRepository.Get(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingResponse(*res))
	}
}
