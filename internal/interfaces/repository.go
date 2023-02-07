package interfaces

import (
	"context"
	"database/sql"
	"mime/multipart"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
)

type IThingRepository interface {
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

type IPlaceRepository interface {
	GetAll(ctx context.Context) ([]models.Place, error)
	Get(ctx context.Context, placeID int) (*models.Place, error)
	GetNestedPlaces(ctx context.Context, placeID int) ([]models.Place, error)
	Add(ctx context.Context, req models.AddPlaceRequest, tx *sql.Tx) (int, error)
	Update(ctx context.Context, req models.UpdatePlaceRequest, tx *sql.Tx) error
	Delete(ctx context.Context, placeID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type IPlaceThingRepository interface {
	GetByThingID(ctx context.Context, thingID int) (*models.PlaceThing, error)
	Add(ctx context.Context, req models.AddPlaceThingRequest, tx *sql.Tx) error
	UpdatePlace(ctx context.Context, req models.UpdatePlaceThingRequest, tx *sql.Tx) error
	DeleteThing(ctx context.Context, thingID int, tx *sql.Tx) error
}

type IPlaceImageRepository interface {
	Add(ctx context.Context, req models.AddPlaceImageRequest, tx *sql.Tx) error
	Get(ctx context.Context, imageID int) (*models.Image, error)
	GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error)
	Delete(ctx context.Context, imageID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type IThingImageRepository interface {
	Add(ctx context.Context, req models.AddThingImageRequest, tx *sql.Tx) error
	Get(ctx context.Context, imageID int) (*models.Image, error)
	GetByThingID(ctx context.Context, thingID int) ([]models.Image, error)
	GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error)
	Delete(ctx context.Context, imageID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type IUserRepository interface {
	Get(ctx context.Context, username string) (*models.User, error)
	//Add(ctx context.Context, req models.AddUserRequest) (int, error)
	//Update(ctx context.Context, req models.UpdateUserRequest) error
	//Delete(ctx context.Context, username string) error
}

type IFileRepository interface {
	Save(fctx *fiber.Ctx, header *multipart.FileHeader, path string) error
	Delete(path string) error
}
