package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.TagRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	tagTableName = "tag"
)

type tagRepository struct {
	db *sql.DB
}

func InitTagRepository(db *sql.DB) interfaces.TagRepository {
	return tagRepository{db: db}
}

func (r tagRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r tagRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r tagRepository) GetAll(ctx context.Context) ([]models.Tag, error) {
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

func (r tagRepository) Get(ctx context.Context, tagID int) (*models.Tag, error) {
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

func (r tagRepository) Add(ctx context.Context, req models.AddTagRequest, tx *sql.Tx) (int, error) {
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

func (r tagRepository) Update(ctx context.Context, req models.UpdateTagRequest, tx *sql.Tx) error {
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

func (r tagRepository) Delete(ctx context.Context, tagID int, tx *sql.Tx) error {
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
