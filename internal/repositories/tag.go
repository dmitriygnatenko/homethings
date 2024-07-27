package repositories

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	tagTableName = "tag"
)

type TagRepository struct {
	db *sql.DB
}

func InitTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r TagRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r TagRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r TagRepository) GetAll(ctx context.Context) ([]models.Tag, error) {
	var res []models.Tag

	query, args, err := sq.Select("id", "title", "style", "created_at", "updated_at").
		From(tagTableName).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resRow := models.Tag{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Title,
			&resRow.Style,
			&resRow.CreatedAt,
			&resRow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, resRow)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r TagRepository) GetByThingID(ctx context.Context, thingID int) ([]models.Tag, error) {
	var res []models.Tag

	query, args, err := sq.Select("t.id", "t.title", "t.style", "t.created_at", "t.updated_at").
		From(tagTableName + " t").
		Join(thingTagTableName + " tt ON tt.tag_id = t.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"tt.thing_id": thingID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resRow := models.Tag{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Title,
			&resRow.Style,
			&resRow.CreatedAt,
			&resRow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, resRow)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r TagRepository) Get(ctx context.Context, tagID int) (*models.Tag, error) {
	query, args, err := sq.Select("id", "title", "style", "created_at", "updated_at").
		From(tagTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": tagID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.Tag

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.ID, &res.Title, &res.Style, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r TagRepository) Add(ctx context.Context, req models.AddTagRequest, tx *sql.Tx) (int, error) {
	query, args, err := sq.Insert(tagTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("title", "style").
		Values(req.Title, req.Style).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, err
	}

	var id int
	if tx == nil {
		err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	} else {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&id)
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r TagRepository) Update(ctx context.Context, req models.UpdateTagRequest, tx *sql.Tx) error {
	query, args, err := sq.Update(tagTableName).
		PlaceholderFormat(sq.Dollar).
		Set("title", req.Title).
		Set("style", req.Style).
		Set("updated_at", "NOW()").
		Where(sq.Eq{"id": req.ID}).
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

func (r TagRepository) Delete(ctx context.Context, tagID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(tagTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": tagID}).
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
