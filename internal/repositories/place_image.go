package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.IPlaceImageRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	placeImageTableName = "place_image"
)

type placeImageRepository struct {
	db *sql.DB
}

func InitPlaceImageRepository(db *sql.DB) interfaces.IPlaceImageRepository {
	return placeImageRepository{db: db}
}

func (r placeImageRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r placeImageRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r placeImageRepository) Add(ctx context.Context, req models.AddPlaceImageRequest, tx *sql.Tx) error {
	query, args, err := sq.Insert(placeImageTableName).
		Columns("place_id", "image").
		Values(req.PlaceID, req.Image).
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

func (r placeImageRepository) Get(ctx context.Context, imageID int) (*models.Image, error) {
	query, args, err := sq.Select("id", "image", "place_id", "created_at").
		From(placeImageTableName).
		Where(sq.Eq{"id": imageID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.Image

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.ID, &res.Image, &res.PlaceID, &res.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r placeImageRepository) GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error) {
	var res []models.Image

	query := "WITH RECURSIVE cte (id, parent_id) AS (" +
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

	rows, err := r.db.QueryContext(ctx, query, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resRow := models.Image{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Image,
			&resRow.PlaceID,
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

func (r placeImageRepository) Delete(ctx context.Context, imageID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(placeImageTableName).
		Where(sq.Eq{"id": imageID}).
		Limit(1).
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
