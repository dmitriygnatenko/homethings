package interfaces

import (
	"context"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

type ITagRepository interface {
	GetAll(ctx context.Context) ([]models.Tag, error)
}
