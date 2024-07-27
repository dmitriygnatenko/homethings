package image

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i FileRepository,ThingImageRepository,PlaceImageRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const fileDateLayout = "2006-01-02-15-04-05"

type (
	FileRepository interface {
		Save(fctx *fiber.Ctx, header *multipart.FileHeader, path string) error
		Delete(path string) error
	}

	ThingImageRepository interface {
		Add(ctx context.Context, req models.AddThingImageRequest, tx *sql.Tx) error
		Get(ctx context.Context, imageID int) (*models.Image, error)
		GetByThingID(ctx context.Context, thingID int) ([]models.Image, error)
		GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error)
		Delete(ctx context.Context, imageID int, tx *sql.Tx) error
		BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
		CommitTx(tx *sql.Tx) error
	}

	PlaceImageRepository interface {
		Add(ctx context.Context, req models.AddPlaceImageRequest, tx *sql.Tx) error
		Get(ctx context.Context, imageID int) (*models.Image, error)
		GetByPlaceID(ctx context.Context, placeID int) ([]models.Image, error)
		Delete(ctx context.Context, imageID int, tx *sql.Tx) error
		BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error)
		CommitTx(tx *sql.Tx) error
	}
)

// @Router 		/api/v1/images [post]
// @Param       place_id formData int false "Place ID"
// @Param       thing_id formData int false "Thing ID"
// @Param       files formData []file true "Files"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add images
// @Tags  		Images
// @security 	APIKey
// @Accept      mpfd
// @Produce     json
func AddImageHandler(
	fileRepository FileRepository,
	thingImageRepository ThingImageRepository,
	placeImageRepository PlaceImageRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		var form *multipart.Form
		var placeID, thingID int
		var files []string
		var err error

		ctx := fctx.Context()

		if form, err = fctx.MultipartForm(); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if formPlace := form.Value["place_id"]; len(formPlace) > 0 {
			placeID, err = strconv.Atoi(formPlace[0])
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
		}

		if formThing := form.Value["thing_id"]; len(formThing) > 0 {
			thingID, err = strconv.Atoi(formThing[0])
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
		}

		date := time.Now().Format(fileDateLayout)
		for _, file := range form.File["files"] {
			filename := "/files/" + date + "_" + helpers.GenerateRandomString(10) + filepath.Ext(file.Filename)

			if err = fileRepository.Save(fctx, file, filename); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			files = append(files, filename)
		}

		if (placeID == 0 && thingID == 0) || len(files) == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		var tx *sql.Tx
		if thingID > 0 {
			tx, err = thingImageRepository.BeginTx(ctx, API.DefaultTxLevel)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			for _, file := range files {
				if err = thingImageRepository.Add(ctx, mappers.ToAddThingImageRequest(thingID, file), tx); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}

			if err = thingImageRepository.CommitTx(tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		} else {
			tx, err = placeImageRepository.BeginTx(ctx, API.DefaultTxLevel)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			for _, file := range files {
				if err = placeImageRepository.Add(ctx, mappers.ToAddPlaceImageRequest(placeID, file), tx); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}

			if err = placeImageRepository.CommitTx(tx); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
