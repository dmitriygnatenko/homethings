package repositories

import (
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

const (
	uploadPath   = "../../web/public"
	resizeWidth  = 1000
	resizeHeight = 800
)

type FileRepository struct{}

func InitFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r FileRepository) Save(fctx *fiber.Ctx, header *multipart.FileHeader, path string) error {
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

func (r FileRepository) Delete(path string) error {
	return os.RemoveAll(getFullPath(path))
}

func getFullPath(path string) string {
	return uploadPath + path
}
