package repositories

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	userTableName = "\"user\""
)

type UserRepository struct {
	db *sql.DB
}

func InitUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) Get(ctx context.Context, username string) (*models.User, error) {
	query, args, err := sq.Select("id", "username", "password", "created_at", "updated_at").
		From(userTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"username": username}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.User
	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.ID, &res.Username, &res.Password, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r UserRepository) Add(ctx context.Context, username string, password string) (int, error) {
	query, args, err := sq.Insert(userTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("username", "password").
		Values(username, password).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, err
	}

	var id int
	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r UserRepository) Update(ctx context.Context, req models.UpdateUserRequest) error {
	qb := sq.Update(userTableName).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", "NOW()").
		Where(sq.Eq{"id": req.ID})

	if req.Username.Valid {
		qb = qb.Set("username", req.Username.String)
	}

	if req.Password.Valid {
		qb = qb.Set("password", req.Password.String)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	return err
}
