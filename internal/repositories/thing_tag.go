package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.ThingTagRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	thingTagTableName = "thing_tag"
)

type thingTagRepository struct {
	db *sql.DB
}

func InitThingTagRepository(db *sql.DB) interfaces.ThingTagRepository {
	return thingTagRepository{db: db}
}

func (r thingTagRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r thingTagRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r thingTagRepository) GetByThingID(ctx context.Context, thingID int) ([]models.Tag, error) {
	var res []models.Tag

	query, args, err := sq.Select("t.id", "t.title", "t.style", "t.created_at", "t.updated_at").
		From(tagTableName + " t").
		Join(thingTagTableName + " tt ON tt.thing_id = t.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"tt.thing_id": thingID}).
		OrderBy("t.updated_at DESC").
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

func (r thingTagRepository) Add(ctx context.Context, req models.AddThingTagRequest, tx *sql.Tx) error {
	query, args, err := sq.Insert(thingTagTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("thing_id", "tag_id").
		Values(req.ThingID, req.TagID).
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

func (r thingTagRepository) Delete(ctx context.Context, req models.DeleteThingTagRequest, tx *sql.Tx) error {
	query, args, err := sq.Delete(thingTagTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"thing_id": req.ThingID}).
		Where(sq.Eq{"tag_id": req.TagID}).
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

func (r thingTagRepository) DeleteByThingID(ctx context.Context, thingID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(thingTagTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"thing_id": thingID}).
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
