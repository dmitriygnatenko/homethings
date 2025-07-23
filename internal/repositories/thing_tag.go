package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/dmitriygnatenko/homethings/internal/models"
)

const thingTagTableName = "thing_tag"

type ThingTagRepository struct {
	db DB
}

func InitThingTagRepository(db DB) *ThingTagRepository {
	return &ThingTagRepository{db: db}
}

func (r ThingTagRepository) Add(ctx context.Context, req models.AddThingTagRequest) error {
	q, v, err := sq.Insert(thingTagTableName).
		PlaceholderFormat(placeholder).
		Columns("thing_id", "tag_id").
		Values(req.ThingID, req.TagID).
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

func (r ThingTagRepository) Delete(ctx context.Context, req models.DeleteThingTagRequest) error {
	q, v, err := sq.Delete(thingTagTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"thing_id": req.ThingID}).
		Where(sq.Eq{"tag_id": req.TagID}).
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

func (r ThingTagRepository) DeleteByTagID(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(thingTagTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"tag_id": id}).
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

func (r ThingTagRepository) DeleteByThingID(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(thingTagTableName).
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

func (r ThingTagRepository) GetByPlaceID(ctx context.Context, id uint64) ([]models.ThingTag, error) {
	var res []models.ThingTag

	q := "WITH RECURSIVE cte (id, parent_id) AS (" +
		"SELECT id, parent_id " +
		"FROM " + placeTableName + " " +
		"WHERE id = ? " +
		"UNION ALL " +
		"SELECT p.id, p.parent_id " +
		"FROM " + placeTableName + " p " +
		"INNER JOIN cte ON p.parent_id = cte.id " +
		")" +
		"SELECT t.id, t.title, t.style, t.created_at, t.updated_at, tt.thing_id " +
		"FROM cte, " + placeThingTableName + " pt, " + thingTagTableName + " tt, " + tagTableName + " t " +
		"WHERE pt.place_id = cte.id AND tt.thing_id = pt.thing_id AND tt.tag_id = t.id " +
		"ORDER BY t.updated_at DESC"

	err := r.db.SelectContext(ctx, &res, q, id)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return res, nil
}
