package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.IThingImageRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	thingImageTableName = "thing_image"
)

type thingImageRepository struct {
	db *sql.DB
}

func InitThingImageRepository(db *sql.DB) interfaces.IThingImageRepository {
	return thingImageRepository{db: db}
}

func (r thingImageRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r thingImageRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r thingImageRepository) Add(ctx context.Context, req models.AddThingImageRequest, tx *sql.Tx) error {
	query, args, err := sq.Insert(thingImageTableName).
		Columns("thing_id", "image").
		Values(req.ThingID, req.Image).
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

func (r thingImageRepository) GetByThingID(ctx context.Context, thingID int) ([]models.Image, error) {
	var res []models.Image

	query, args, err := sq.Select("id", "image", "thing_id", "created_at").
		From(thingImageTableName).
		Where(sq.Eq{"thing_id": thingID}).
		OrderBy("created_at desc").
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
		resRow := models.Image{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Image,
			&resRow.ThingID,
			&resRow.CreatedAt,
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
