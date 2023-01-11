package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.IFileRepository -o ./mocks/ -s "_minimock.go"

import (
	"mime/multipart"
	"os"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

type fileRepository struct{}

func InitFileRepository() interfaces.IFileRepository {
	return fileRepository{}
}

func (r fileRepository) Save(fctx *fiber.Ctx, header *multipart.FileHeader, path string) error {
	return fctx.SaveFile(header, path)
}

func (r fileRepository) Delete(path string) error {
	return os.RemoveAll(path)
}
