package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.IPlaceRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	placeTableName = "place"
)

type placeRepository struct {
	db *sql.DB
}

func InitPlaceRepository(db *sql.DB) interfaces.IPlaceRepository {
	return placeRepository{db: db}
}

func (r placeRepository) GetAll(ctx context.Context) ([]models.Place, error) {
	var res []models.Place

	query, args, err := sq.Select("id", "parent_id", "title", "created_at", "updated_at").
		From(placeTableName).
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
		resRow := models.Place{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.ParentID,
			&resRow.Title,
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

func (r placeRepository) Get(ctx context.Context, placeID int) (*models.Place, error) {
	query, args, err := sq.Select("id", "parent_id", "title", "created_at", "updated_at").
		From(placeTableName).
		Where(sq.Eq{"id": placeID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.Place

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.ID, &res.ParentID, &res.Title, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r placeRepository) Add(ctx context.Context, req models.AddPlaceRequest, tx *sql.Tx) (int, error) {
	query, args, err := sq.Insert(placeTableName).
		Columns("title", "parent_id").
		Values(req.Title, req.ParentID).
		ToSql()

	if err != nil {
		return 0, err
	}

	var res sql.Result

	if tx == nil {
		res, err = r.db.ExecContext(ctx, query, args...)
	} else {
		res, err = tx.ExecContext(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r placeRepository) Update(ctx context.Context, req models.UpdatePlaceRequest, tx *sql.Tx) error {
	query, args, err := sq.Update(placeTableName).
		Set("title", req.Title).
		Set("parent_id", req.ParentID).
		Where(sq.Eq{"id": req.ID}).ToSql()

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
