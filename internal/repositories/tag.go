package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const tagTableName = "tag"

var tagTableFields = []string{"id", "title", "style", "created_at", "updated_at"}

type TagRepository struct {
	db DB
}

func InitTagRepository(db DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r TagRepository) GetAll(ctx context.Context) ([]models.Tag, error) {
	var res []models.Tag

	q, v, err := sq.Select(tagTableFields...).
		From(tagTableName).
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

func (r TagRepository) GetByThingID(ctx context.Context, id uint64) ([]models.Tag, error) {
	var res []models.Tag

	q, v, err := sq.Select("t.id", "t.title", "t.style", "t.created_at", "t.updated_at").
		From(tagTableName + " t").
		Join(thingTagTableName + " tt ON tt.tag_id = t.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"tt.thing_id": id}).
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

func (r TagRepository) Get(ctx context.Context, id uint64) (*models.Tag, error) {
	q, v, err := sq.Select(tagTableFields...).
		From(tagTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res models.Tag

	err = r.db.GetContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &res, nil
}

func (r TagRepository) Add(ctx context.Context, req models.AddTagRequest) (uint64, error) {
	q, v, err := sq.Insert(tagTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("title", "style").
		Values(req.Title, req.Style).
		Suffix("RETURNING id").
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
		return 0, fmt.Errorf("last insert ID: %w", err)
	}

	return uint64(id), nil
}

func (r TagRepository) Update(ctx context.Context, req models.UpdateTagRequest) error {
	q, v, err := sq.Update(tagTableName).
		PlaceholderFormat(sq.Dollar).
		Set("title", req.Title).
		Set("style", req.Style).
		Set("updated_at", "NOW()").
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

func (r TagRepository) Delete(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(tagTableName).
		PlaceholderFormat(sq.Dollar).
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
