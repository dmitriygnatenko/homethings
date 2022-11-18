package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.IThingRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
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

func (r thingRepository) Get(ctx context.Context, id int) (*models.Thing, error) {
	var res models.Thing

	query := "SELECT id, title, description, created_at, updated_at FROM " + thingTableName + " WHERE id = ?"

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&res.ID, &res.Title, &res.Description, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r thingRepository) Add(ctx context.Context, req models.AddThingRequest, tx *sql.Tx) (int, error) {
	var res sql.Result
	var err error

	query := "INSERT INTO " + thingTableName + " (title, description) VALUES (?, ?)"

	if tx == nil {
		res, err = r.db.ExecContext(ctx, query, req.Title, req.Description)
	} else {
		res, err = tx.ExecContext(ctx, query, req.Title, req.Description)
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
	return nil
}

func (r thingRepository) Delete(ctx context.Context, id int, tx *sql.Tx) error {
	return nil
}
