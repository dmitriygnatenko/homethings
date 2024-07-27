package repositories

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	placeTableName = "place"
)

type PlaceRepository struct {
	db *sql.DB
}

func InitPlaceRepository(db *sql.DB) *PlaceRepository {
	return &PlaceRepository{db: db}
}

func (r PlaceRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r PlaceRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r PlaceRepository) GetAll(ctx context.Context) ([]models.Place, error) {
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

func (r PlaceRepository) GetNestedPlaces(ctx context.Context, placeID int) ([]models.Place, error) {
	var res []models.Place

	query, args, err := sq.Select("id", "parent_id", "title", "created_at", "updated_at").
		From(placeTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"parent_id": placeID}).
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

func (r PlaceRepository) Get(ctx context.Context, placeID int) (*models.Place, error) {
	query, args, err := sq.Select("id", "parent_id", "title", "created_at", "updated_at").
		From(placeTableName).
		PlaceholderFormat(sq.Dollar).
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

func (r PlaceRepository) Add(ctx context.Context, req models.AddPlaceRequest, tx *sql.Tx) (int, error) {
	query, args, err := sq.Insert(placeTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("title", "parent_id").
		Values(req.Title, req.ParentID).
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

func (r PlaceRepository) Update(ctx context.Context, req models.UpdatePlaceRequest, tx *sql.Tx) error {
	query, args, err := sq.Update(placeTableName).
		PlaceholderFormat(sq.Dollar).
		Set("title", req.Title).
		Set("parent_id", req.ParentID).
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

func (r PlaceRepository) Delete(ctx context.Context, placeID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(placeTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": placeID}).
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
