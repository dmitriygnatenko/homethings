package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.ITagRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	tagTableName = "tag"
)

type tagRepository struct {
	db *sql.DB
}

func InitTagRepository(db *sql.DB) interfaces.ITagRepository {
	return tagRepository{db: db}
}

func (t tagRepository) GetAll(ctx context.Context) ([]models.Tag, error) {
	var res []models.Tag

	query := "SELECT id, title, created_at, updated_at FROM " + tagTableName

	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row := models.Tag{}

		err = rows.Scan(
			&row.ID,
			&row.Title,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, row)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}
