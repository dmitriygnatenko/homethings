package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const userTableName = "\"user\""

var userTableFields = []string{"id", "username", "password", "created_at", "updated_at"}

type UserRepository struct {
	db DB
}

func InitUserRepository(db DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) Get(ctx context.Context, username string) (*models.User, error) {
	q, v, err := sq.Select(userTableFields...).
		From(userTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"username": username}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res models.User

	err = r.db.GetContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &res, nil
}

func (r UserRepository) Add(ctx context.Context, username string, password string) (uint64, error) {
	q, v, err := sq.Insert(userTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("username", "password").
		Values(username, password).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	var id uint64
	err = r.db.QueryRowContext(ctx, q, v...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
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

	q, v, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
