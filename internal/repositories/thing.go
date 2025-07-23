package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/dmitriygnatenko/homethings/internal/models"
)

const thingTableName = "thing"

type ThingRepository struct {
	db DB
}

func InitThingRepository(db DB) *ThingRepository {
	return &ThingRepository{db: db}
}

func (r ThingRepository) Get(ctx context.Context, id uint64) (*models.Thing, error) {
	q, v, err := sq.Select("t.id", "t.title", "t.description", "t.created_at", "t.updated_at", "p.place_id").
		From(thingTableName + " t").
		Join(placeThingTableName + " p ON p.thing_id = t.id").
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res models.Thing

	err = r.db.GetContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &res, nil
}

func (r ThingRepository) Search(ctx context.Context, search string) ([]models.Thing, error) {
	var res []models.Thing

	s := fmt.Sprint("%", search, "%")

	q, v, err := sq.Select("t.id", "t.title", "t.description", "t.created_at", "t.updated_at", "p.place_id").
		From(thingTableName+" t").
		Join(placeThingTableName+" p ON p.thing_id = t.id").
		PlaceholderFormat(placeholder).
		Where("LOWER(t.title) LIKE LOWER(?) OR LOWER(t.description) LIKE LOWER(?)", s, s).
		OrderBy("t.updated_at DESC").
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

func (r ThingRepository) GetByPlaceID(ctx context.Context, id uint64) ([]models.Thing, error) {
	var res []models.Thing

	q, v, err := sq.Select("t.id", "t.title", "t.description", "t.created_at", "t.updated_at", "p.place_id").
		From(thingTableName + " t").
		Join(placeThingTableName + " p ON p.thing_id = t.id").
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"p.place_id": id}).
		OrderBy("t.updated_at DESC").
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

// GetAllByPlaceID return things by place ID and all child places.
func (r ThingRepository) GetAllByPlaceID(ctx context.Context, id uint64) ([]models.Thing, error) {
	var res []models.Thing

	q := "WITH RECURSIVE cte (id, parent_id) AS (" +
		"SELECT id, parent_id " +
		"FROM " + placeTableName + " " +
		"WHERE id = ? " +
		"UNION ALL " +
		"SELECT p.id, p.parent_id " +
		"FROM " + placeTableName + " p " +
		"INNER JOIN cte ON p.parent_id = cte.id " +
		")" +
		"SELECT t.id, t.title, t.description, t.created_at, t.updated_at, pt.place_id " +
		"FROM cte, " + placeThingTableName + " pt, " + thingTableName + " t " +
		"WHERE pt.place_id = cte.id and t.id = pt.thing_id " +
		"ORDER BY t.updated_at DESC"

	err := r.db.SelectContext(ctx, &res, q, id)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return res, nil
}

func (r ThingRepository) Add(ctx context.Context, req models.AddThingRequest) (uint64, error) {
	q, v, err := sq.Insert(thingTableName).
		PlaceholderFormat(placeholder).
		Columns("title", "description").
		Values(req.Title, req.Description).
		ToSql()

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

func (r ThingRepository) Update(ctx context.Context, req models.UpdateThingRequest) error {
	q, v, err := sq.Update(thingTableName).
		PlaceholderFormat(placeholder).
		Set("title", req.Title).
		Set("description", req.Description).
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

func (r ThingRepository) Delete(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(thingTableName).
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
