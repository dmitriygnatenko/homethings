package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

const thingImageTableName = "thing_image"

type ThingImageRepository struct {
	db DB
}

func InitThingImageRepository(db DB) *ThingImageRepository {
	return &ThingImageRepository{db: db}
}

func (r ThingImageRepository) Add(ctx context.Context, req models.AddThingImageRequest) error {
	q, v, err := sq.Insert(thingImageTableName).
		PlaceholderFormat(placeholder).
		Columns("thing_id", "image").
		Values(req.ThingID, req.Image).
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

func (r ThingImageRepository) Get(ctx context.Context, id uint64) (*models.Image, error) {
	q, v, err := sq.Select("id", "image", "thing_id", "created_at").
		From(thingImageTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res models.Image

	err = r.db.GetContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &res, nil
}

func (r ThingImageRepository) GetByThingID(ctx context.Context, id uint64) ([]models.Image, error) {
	var res []models.Image

	q, v, err := sq.Select("id", "image", "thing_id", "created_at").
		From(thingImageTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"thing_id": id}).
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	err = r.db.SelectContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return res, nil
}

func (r ThingImageRepository) GetByPlaceID(ctx context.Context, id uint64) ([]models.Image, error) {
	var res []models.Image

	q := "WITH RECURSIVE cte (id, parent_id) AS (" +
		"SELECT id, parent_id " +
		"FROM " + placeTableName + " " +
		"WHERE id = ? " +
		"UNION ALL " +
		"SELECT p.id, p.parent_id " +
		"FROM " + placeTableName + " p " +
		"INNER JOIN cte ON p.parent_id = cte.id " +
		")" +
		"SELECT ti.id, ti.image, ti.thing_id, ti.created_at " +
		"FROM cte, " + placeThingTableName + " pt, " + thingImageTableName + " ti " +
		"WHERE pt.place_id = cte.id AND pt.thing_id = ti.thing_id " +
		"ORDER BY ti.created_at DESC"

	err := r.db.SelectContext(ctx, &res, q, id)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return res, nil
}

func (r ThingImageRepository) Delete(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(thingImageTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"id": id}).
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
