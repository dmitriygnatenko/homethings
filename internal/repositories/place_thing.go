package repositories

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	placeThingTableName = "place_thing"
)

type PlaceThingRepository struct {
	db *sql.DB
}

func InitPlaceThingRepository(db *sql.DB) *PlaceThingRepository {
	return &PlaceThingRepository{db: db}
}

func (r PlaceThingRepository) GetByThingID(ctx context.Context, thingID int) (*models.PlaceThing, error) {
	query, args, err := sq.Select("place_id", "thing_id", "created_at").
		From(placeThingTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"thing_id": thingID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.PlaceThing

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.PlaceID, &res.ThingID, &res.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r PlaceThingRepository) Add(ctx context.Context, req models.AddPlaceThingRequest, tx *sql.Tx) error {
	query, args, err := sq.Insert(placeThingTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("place_id", "thing_id").
		Values(req.PlaceID, req.ThingID).
		ToSql()

	if err != nil {
		return err
	}

	if tx == nil {
		_, err = r.db.ExecContext(ctx, query, args...)
	} else {
		_, err = tx.ExecContext(ctx, query, args...)
	}

	return err
}

func (r PlaceThingRepository) UpdatePlace(ctx context.Context, req models.UpdatePlaceThingRequest, tx *sql.Tx) error {
	query, args, err := sq.Update(placeThingTableName).
		PlaceholderFormat(sq.Dollar).
		Set("place_id", req.PlaceID).
		Set("updated_at", "NOW()").
		Where(sq.Eq{"thing_id": req.ThingID}).
		ToSql()

	if err != nil {
		return err
	}

	if tx == nil {
		_, err = r.db.ExecContext(ctx, query, args...)
	} else {
		_, err = tx.ExecContext(ctx, query, args...)
	}

	return err
}

func (r PlaceThingRepository) DeleteThing(ctx context.Context, id int, tx *sql.Tx) error {
	query, args, err := sq.Delete(placeThingTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"thing_id": id}).
		ToSql()

	if err != nil {
		return err
	}

	if tx == nil {
		_, err = r.db.ExecContext(ctx, query, args...)
	} else {
		_, err = tx.ExecContext(ctx, query, args...)
	}

	return err
}
