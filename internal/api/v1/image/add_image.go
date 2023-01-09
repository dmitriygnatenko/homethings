package image

import (
	"database/sql"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

const (
	uploadPath     = "../../web"
	fileDateLayout = "2006-01-02-15-04-05"
)

// @Router 		/v1/images [post]
// @Param       place_id formData int false "Place ID"
// @Param       thing_id formData int false "Thing ID"
// @Param       files formData []file true "Files"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add images
// @Tags  		Images
// @security 	BasicAuth
// @Accept      mpfd
// @Produce     json
func AddImageHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		var form *multipart.Form
		var placeID, thingID int
		var files []string
		var err error

		ctx := fctx.Context()

		if form, err = fctx.MultipartForm(); err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		if formPlace := form.Value["place_id"]; len(formPlace) > 0 {
			placeID, err = strconv.Atoi(formPlace[0])
			if err != nil {
				return factory.CreateBadRequestResponse(fctx, err)
			}
		}

		if formThing := form.Value["thing_id"]; len(formThing) > 0 {
			thingID, err = strconv.Atoi(formThing[0])
			if err != nil {
				return factory.CreateBadRequestResponse(fctx, err)
			}
		}

		date := time.Now().Format(fileDateLayout)
		for _, file := range form.File["files"] {
			filename := "/files/" + date + "_" + helpers.GenerateRandomString(10) + filepath.Ext(file.Filename)

			if err = sp.GetFileRepository().Save(fctx, file, uploadPath+filename); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}

			files = append(files, filename)
		}

		if (placeID == 0 && thingID == 0) || len(files) == 0 {
			return factory.CreateBadRequestResponse(fctx, nil)
		}

		var tx *sql.Tx
		if thingID > 0 {
			tx, err = sp.GetThingImageRepository().BeginTx(ctx, API.DefaultTxLevel)
			if err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}

			for _, file := range files {
				if err = sp.GetThingImageRepository().Add(ctx, mappers.ConvertToAddThingImageRequestModel(thingID, file), tx); err != nil {
					return factory.CreateInternalErrorResponse(fctx, err)
				}
			}

			if err = sp.GetThingImageRepository().CommitTx(tx); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		} else {
			tx, err = sp.GetPlaceImageRepository().BeginTx(ctx, API.DefaultTxLevel)
			if err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}

			for _, file := range files {
				if err = sp.GetPlaceImageRepository().Add(ctx, mappers.ConvertToAddPlaceImageRequestModel(placeID, file), tx); err != nil {
					return factory.CreateInternalErrorResponse(fctx, err)
				}
			}

			if err = sp.GetPlaceImageRepository().CommitTx(tx); err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
		}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
