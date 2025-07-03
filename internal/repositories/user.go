package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const userTableName = "user"

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
		PlaceholderFormat(placeholder).
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
		PlaceholderFormat(placeholder).
		Columns("username", "password").
		Values(username, password).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return uint64(id), nil
}

func (r UserRepository) Update(ctx context.Context, req models.UpdateUserRequest) error {
	qb := sq.Update(userTableName).
		PlaceholderFormat(placeholder).
		Set("updated_at", time.Now()).
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
