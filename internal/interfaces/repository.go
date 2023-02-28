package interfaces

import (
	"context"
	"database/sql"
	"mime/multipart"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
)

type ThingRepository interface {
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

type PlaceRepository interface {
	GetAll(ctx context.Context) ([]models.Place, error)
	Get(ctx context.Context, placeID int) (*models.Place, error)
	GetNestedPlaces(ctx context.Context, placeID int) ([]models.Place, error)
	Add(ctx context.Context, req models.AddPlaceRequest, tx *sql.Tx) (int, error)
	Update(ctx context.Context, req models.UpdatePlaceRequest, tx *sql.Tx) error
	Delete(ctx context.Context, placeID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type PlaceThingRepository interface {
	GetByThingID(ctx context.Context, thingID int) (*models.PlaceThing, error)
	Add(ctx context.Context, req models.AddPlaceThingRequest, tx *sql.Tx) error
	UpdatePlace(ctx context.Context, req models.UpdatePlaceThingRequest, tx *sql.Tx) error
	DeleteThing(ctx context.Context, thingID int, tx *sql.Tx) error
}

type PlaceImageRepository interface {
	Add(ctx context.Context, req models.AddPlaceImageRequest, tx *sql.Tx) error
	Get(ctx context.Context, imageID int) (*models.Image, error)
	GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error)
	Delete(ctx context.Context, imageID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type ThingImageRepository interface {
	Add(ctx context.Context, req models.AddThingImageRequest, tx *sql.Tx) error
	Get(ctx context.Context, imageID int) (*models.Image, error)
	GetByThingID(ctx context.Context, thingID int) ([]models.Image, error)
	GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error)
	Delete(ctx context.Context, imageID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type UserRepository interface {
	Get(ctx context.Context, username string) (*models.User, error)
	Add(ctx context.Context, username string, password string) (int, error)
	Update(ctx context.Context, req models.UpdateUserRequest) error
}

type TagRepository interface {
	GetAll(ctx context.Context) ([]models.Tag, error)
	Get(ctx context.Context, tagID int) (*models.Tag, error)
	Add(ctx context.Context, req models.AddTagRequest, tx *sql.Tx) (int, error)
	Update(ctx context.Context, req models.UpdateTagRequest, tx *sql.Tx) error
	Delete(ctx context.Context, tagID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type ThingTagRepository interface {
	GetByThingID(ctx context.Context, thingID int) ([]models.ThingTag, error)
	GetByPlaceID(ctx context.Context, placeID int) ([]models.ThingTag, error)
	Add(ctx context.Context, req models.AddThingTagRequest, tx *sql.Tx) error
	Delete(ctx context.Context, req models.DeleteThingTagRequest, tx *sql.Tx) error
	DeleteByThingID(ctx context.Context, thingID int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
}

type FileRepository interface {
	Save(fctx *fiber.Ctx, header *multipart.FileHeader, path string) error
	Delete(path string) error
}
