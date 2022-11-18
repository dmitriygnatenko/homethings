package interfaces

import (
	"context"
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

type IThingRepository interface {
	Get(ctx context.Context, id int) (*models.Thing, error)
	Add(ctx context.Context, req models.AddThingRequest, tx *sql.Tx) (int, error)
	Update(ctx context.Context, req models.UpdateThingRequest, tx *sql.Tx) error
	Delete(ctx context.Context, id int, tx *sql.Tx) error
	BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
}

type IPlaceThingRepository interface {
	Add(ctx context.Context, req models.AddPlaceThingRequest, tx *sql.Tx) error
}

type ITagRepository interface {
	GetAll(ctx context.Context) ([]models.Tag, error)
}
