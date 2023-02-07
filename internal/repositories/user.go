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
	userTableName = "\"user\""
)

type userRepository struct {
	db *sql.DB
}

func InitUserRepository(db *sql.DB) interfaces.IUserRepository {
	return userRepository{db: db}
}

func (r userRepository) Get(ctx context.Context, username string) (*models.User, error) {
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

func (r userRepository) Add(ctx context.Context, username string, password string) (int, error) {
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

//func (r userRepository) Update(ctx context.Context, req models.UpdateUserRequest) error {
//	query, args, err := sq.Update(userTableName).
//		PlaceholderFormat(sq.Dollar).
//		Set("password", req.Password).
//		Set("updated_at", "NOW()").
//		Where(sq.Eq{"username": req.Username}).
//		ToSql()
//
//	if err != nil {
//		return err
//	}
//
//	_, err = r.db.ExecContext(ctx, query, args...)
//
//	return err
//}
