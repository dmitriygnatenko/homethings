package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.IThingRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	thingTableName = "thing"
)

type thingRepository struct {
	db *sql.DB
}

func InitThingRepository(db *sql.DB) interfaces.IThingRepository {
	return thingRepository{db: db}
}

func (r thingRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r thingRepository) Get(ctx context.Context, thingID int) (*models.Thing, error) {
	query, args, err := sq.Select("thingID", "title", "description", "created_at", "updated_at").
		From(thingTableName).
		Where(sq.Eq{"id": thingID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.Thing

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r thingRepository) GetByPlaceID(ctx context.Context, placeID int) ([]models.Thing, error) {
	var res []models.Thing

	query, args, err := sq.Select("t.id", "t.title", "t.description", "t.created_at", "t.updated_at").
		From(thingTableName + " t").
		Join(placeThingTableName + " p ON p.thing_id = t.id").
		Where(sq.Eq{"p.place_id": placeID}).
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
		resRow := models.Thing{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Title,
			&resRow.Description,
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

func (r thingRepository) Add(ctx context.Context, req models.AddThingRequest, tx *sql.Tx) (int, error) {
	query, args, err := sq.Insert(thingTableName).
		Columns("title", "description").
		Values(req.Title, req.Description).
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

func (r thingRepository) Update(ctx context.Context, req models.UpdateThingRequest, tx *sql.Tx) error {
	builder := sq.Update(thingTableName)

	if req.Title.Valid {
		builder = builder.Set("title", req.Title.String)
	}

	if req.Description.Valid {
		builder = builder.Set("description", req.Description.String)
	}

	query, args, err := builder.Where(sq.Eq{"id": req.ID}).ToSql()
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

func (r thingRepository) Delete(ctx context.Context, thingID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(thingTableName).
		Where(sq.Eq{"id": thingID}).
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
