package repositories

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const placeTableName = "place"

var placeTableFields = []string{"id", "parent_id", "title", "created_at", "updated_at"}

type PlaceRepository struct {
	db DB
}

func InitPlaceRepository(db DB) *PlaceRepository {
	return &PlaceRepository{db: db}
}

func (r PlaceRepository) GetAll(ctx context.Context) ([]models.Place, error) {
	var res []models.Place

	q, v, err := sq.Select(placeTableFields...).
		From(placeTableName).
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

func (r PlaceRepository) GetNestedPlaces(ctx context.Context, id uint64) ([]models.Place, error) {
	var res []models.Place

	q, v, err := sq.Select(placeTableFields...).
		From(placeTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"parent_id": id}).
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

func (r PlaceRepository) Get(ctx context.Context, id uint64) (*models.Place, error) {
	q, v, err := sq.Select(placeTableFields...).
		From(placeTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res models.Place

	err = r.db.GetContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &res, nil
}

func (r PlaceRepository) Add(ctx context.Context, req models.AddPlaceRequest) (uint64, error) {
	builder := sq.Insert(placeTableName).
		PlaceholderFormat(placeholder).
		Columns("title", "parent_id").
		Values(req.Title, req.ParentID)

	q, v, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return uint64(id), nil
}

func (r PlaceRepository) Update(ctx context.Context, req models.UpdatePlaceRequest) error {
	q, v, err := sq.Update(placeTableName).
		PlaceholderFormat(placeholder).
		Set("title", req.Title).
		Set("parent_id", req.ParentID).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": req.ID}).
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

func (r PlaceRepository) Delete(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(placeTableName).
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
