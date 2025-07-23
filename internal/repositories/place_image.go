package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/dmitriygnatenko/homethings/internal/models"
)

const placeImageTableName = "place_image"

type PlaceImageRepository struct {
	db DB
}

func InitPlaceImageRepository(db DB) *PlaceImageRepository {
	return &PlaceImageRepository{db: db}
}

func (r PlaceImageRepository) Add(ctx context.Context, req models.AddPlaceImageRequest) error {
	q, v, err := sq.Insert(placeImageTableName).
		PlaceholderFormat(placeholder).
		Columns("place_id", "image").
		Values(req.PlaceID, req.Image).
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

func (r PlaceImageRepository) Get(ctx context.Context, id uint64) (*models.Image, error) {
	q, v, err := sq.Select("id", "image", "place_id", "created_at").
		From(placeImageTableName).
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

func (r PlaceImageRepository) GetByPlaceID(ctx context.Context, id uint64) ([]models.Image, error) {
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
		"SELECT pi.id, pi.image, pi.place_id, pi.created_at " +
		"FROM cte, " + placeImageTableName + " pi " +
		"WHERE pi.place_id = cte.id " +
		"ORDER BY pi.created_at DESC"

	err := r.db.SelectContext(ctx, &res, q, id)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return res, nil
}

func (r PlaceImageRepository) Delete(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(placeImageTableName).
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
