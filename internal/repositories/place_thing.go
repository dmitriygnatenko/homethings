package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

const placeThingTableName = "place_thing"

type PlaceThingRepository struct {
	db DB
}

func InitPlaceThingRepository(db DB) *PlaceThingRepository {
	return &PlaceThingRepository{db: db}
}

func (r PlaceThingRepository) GetByThingID(ctx context.Context, id uint64) (*models.PlaceThing, error) {
	q, v, err := sq.Select("place_id", "thing_id", "created_at").
		From(placeThingTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"thing_id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res models.PlaceThing

	err = r.db.GetContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &res, nil
}

func (r PlaceThingRepository) Add(ctx context.Context, req models.AddPlaceThingRequest) error {
	q, v, err := sq.Insert(placeThingTableName).
		PlaceholderFormat(placeholder).
		Columns("place_id", "thing_id").
		Values(req.PlaceID, req.ThingID).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r PlaceThingRepository) UpdatePlace(ctx context.Context, req models.UpdatePlaceThingRequest) error {
	q, v, err := sq.Update(placeThingTableName).
		PlaceholderFormat(placeholder).
		Set("place_id", req.PlaceID).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"thing_id": req.ThingID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r PlaceThingRepository) DeleteThing(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(placeThingTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"thing_id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
