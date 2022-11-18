package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.IPlaceThingRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	placeThingTableName = "place_thing"
)

type placeThingRepository struct {
	db *sql.DB
}

func InitPlaceThingRepository(db *sql.DB) interfaces.IPlaceThingRepository {
	return placeThingRepository{db: db}
}

func (r placeThingRepository) Add(ctx context.Context, req models.AddPlaceThingRequest, tx *sql.Tx) error {
	var err error

	query := "INSERT INTO " + placeThingTableName + " (place_id, thing_id) VALUES (?, ?)"

	if tx == nil {
		_, err = r.db.ExecContext(ctx, query, req.PlaceID, req.ThingID)
	} else {
		_, err = tx.ExecContext(ctx, query, req.PlaceID, req.ThingID)
	}

	return err
}
