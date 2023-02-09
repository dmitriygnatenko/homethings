package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.ThingRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func InitThingRepository(db *sql.DB) interfaces.ThingRepository {
	return thingRepository{db: db}
}

func (r thingRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r thingRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r thingRepository) Get(ctx context.Context, thingID int) (*models.Thing, error) {
	query, args, err := sq.Select("t.id", "t.title", "t.description", "t.created_at", "t.updated_at", "p.place_id").
		From(thingTableName + " t").
		Join(placeThingTableName + " p ON p.thing_id = t.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": thingID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.Thing

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.CreatedAt, &res.UpdatedAt, &res.PlaceID)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r thingRepository) Search(ctx context.Context, search string) ([]models.Thing, error) {
	var res []models.Thing

	query, args, err := sq.Select("t.id", "t.title", "t.description", "t.created_at", "t.updated_at", "p.place_id").
		From(thingTableName+" t").
		Join(placeThingTableName+" p ON p.thing_id = t.id").
		PlaceholderFormat(sq.Dollar).
		Where("t.title ILIKE ?", fmt.Sprint("%", search, "%")).
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
		resRow := models.Thing{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Title,
			&resRow.Description,
			&resRow.CreatedAt,
			&resRow.UpdatedAt,
			&resRow.PlaceID,
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

func (r thingRepository) GetByPlaceID(ctx context.Context, placeID int) ([]models.Thing, error) {
	var res []models.Thing

	query, args, err := sq.Select("t.id", "t.title", "t.description", "t.created_at", "t.updated_at", "p.place_id").
		From(thingTableName + " t").
		Join(placeThingTableName + " p ON p.thing_id = t.id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"p.place_id": placeID}).
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
		resRow := models.Thing{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Title,
			&resRow.Description,
			&resRow.CreatedAt,
			&resRow.UpdatedAt,
			&resRow.PlaceID,
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

// GetAllByPlaceID return things by place ID and all child places
func (r thingRepository) GetAllByPlaceID(ctx context.Context, placeID int) ([]models.Thing, error) {
	var res []models.Thing

	query := "WITH RECURSIVE cte (id, parent_id) AS (" +
		"SELECT id, parent_id " +
		"FROM " + placeTableName + " " +
		"WHERE id = $1 " +
		"UNION ALL " +
		"SELECT p.id, p.parent_id " +
		"FROM " + placeTableName + " p " +
		"INNER JOIN cte ON p.parent_id = cte.id " +
		")" +
		"SELECT t.id, t.title, t.description, t.created_at, t.updated_at, pt.place_id " +
		"FROM cte, " + placeThingTableName + " pt, " + thingTableName + " t " +
		"WHERE pt.place_id = cte.id and t.id = pt.thing_id " +
		"ORDER BY t.updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query, placeID)
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
			&resRow.PlaceID,
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
		PlaceholderFormat(sq.Dollar).
		Columns("title", "description").
		Values(req.Title, req.Description).
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

func (r thingRepository) Update(ctx context.Context, req models.UpdateThingRequest, tx *sql.Tx) error {
	query, args, err := sq.Update(thingTableName).
		PlaceholderFormat(sq.Dollar).
		Set("title", req.Title).
		Set("description", req.Description).
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

func (r thingRepository) Delete(ctx context.Context, thingID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(thingTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": thingID}).
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
