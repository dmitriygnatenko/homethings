package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.FileRepository -o ./mocks/ -s "_minimock.go"

import (
	"mime/multipart"
	"os"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

const (
	uploadPath   = "../../web/public"
	resizeWidth  = 1000
	resizeHeight = 800
)

type fileRepository struct{}

func InitFileRepository() interfaces.FileRepository {
	return fileRepository{}
}

func (r fileRepository) Save(fctx *fiber.Ctx, header *multipart.FileHeader, path string) error {
	fullPath := getFullPath(path)

	if err := fctx.SaveFile(header, fullPath); err != nil {
		return err
	}

	// Resize image
	src, err := imaging.Open(fullPath)
	if err != nil {
		return err
	}

	origHeight := src.Bounds().Size().Y
	origWidth := src.Bounds().Size().X
	var newHeight, newWidth int

	if origWidth > origHeight {
		newWidth = resizeWidth
	} else {
		newHeight = resizeHeight
	}

	dst := imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)
	if err = imaging.Save(dst, fullPath); err != nil {
		return err
	}

	return nil
}

func (r fileRepository) Delete(path string) error {
	return os.RemoveAll(getFullPath(path))
}

func getFullPath(path string) string {
	return uploadPath + path
}
