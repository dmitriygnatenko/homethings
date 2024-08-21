package image

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i TransactionManager,FileRepository,ThingImageRepository,PlaceImageRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/random"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const fileDateLayout = "2006-01-02-15-04-05"

type (
	TransactionManager interface {
		ReadCommitted(context.Context, func(ctx context.Context) error) error
	}

	FileRepository interface {
		Save(fctx *fiber.Ctx, header *multipart.FileHeader, path string) error
		Delete(path string) error
	}

	ThingImageRepository interface {
		Add(ctx context.Context, req models.AddThingImageRequest) error
		Get(ctx context.Context, id uint64) (*models.Image, error)
		GetByThingID(ctx context.Context, id uint64) ([]models.Image, error)
		GetByPlaceID(ctx context.Context, id uint64) ([]models.Image, error)
		Delete(ctx context.Context, id uint64) error
	}

	PlaceImageRepository interface {
		Add(ctx context.Context, req models.AddPlaceImageRequest) error
		Get(ctx context.Context, id uint64) (*models.Image, error)
		GetByPlaceID(ctx context.Context, id uint64) ([]models.Image, error)
		Delete(ctx context.Context, id uint64) error
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
	tm TransactionManager,
	fileRepository FileRepository,
	thingImageRepository ThingImageRepository,
	placeImageRepository PlaceImageRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		var form *multipart.Form
		var placeID, thingID uint64
		var files []string
		var err error

		ctx := fctx.Context()

		if form, err = fctx.MultipartForm(); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if formPlace := form.Value["place_id"]; len(formPlace) > 0 {
			placeID, err = strconv.ParseUint(formPlace[0], 0, 64)
			if err != nil {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
		}

		if formThing := form.Value["thing_id"]; len(formThing) > 0 {
			thingID, err = strconv.ParseUint(formThing[0], 0, 64)
			if err != nil {
				logger.Info(ctx, err.Error())
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
		}

		date := time.Now().Format(fileDateLayout)
		for _, file := range form.File["files"] {
			filename := "/files/" + date + "_" + random.GenerateRandomString(10) + filepath.Ext(file.Filename)

			if err = fileRepository.Save(fctx, file, filename); err != nil {
				logger.Error(ctx, err.Error())
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			files = append(files, filename)
		}

		if (placeID == 0 && thingID == 0) || len(files) == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		err = tm.ReadCommitted(ctx, func(ctx context.Context) error {
			if thingID > 0 {
				for _, file := range files {
					txErr := thingImageRepository.Add(ctx, mappers.ToAddThingImageRequest(thingID, file))
					if txErr != nil {
						return txErr
					}
				}
			} else {
				for _, file := range files {
					txErr := placeImageRepository.Add(ctx, mappers.ToAddPlaceImageRequest(placeID, file))
					if txErr != nil {
						return txErr
					}
				}
			}

			return nil
		})

		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
