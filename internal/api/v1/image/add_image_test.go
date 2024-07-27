package image

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/image/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

// nolint:errcheck
func TestAddImageHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        []byte
	}

	var (
		placeID   = gofakeit.Number(1, 1000)
		thingID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())
	)

	// Correct request for adding place image
	addPlaceCorrectBody := &bytes.Buffer{}
	addPlaceCorrectWriter := multipart.NewWriter(addPlaceCorrectBody)
	addPlaceCorrectWriter.WriteField("place_id", strconv.Itoa(placeID))
	addPlaceCorrectWriter.CreateFormFile("files", gofakeit.Word())
	addPlaceCorrectContentType := addPlaceCorrectWriter.FormDataContentType()
	addPlaceCorrectWriter.Close()

	// Correct request for adding thing image
	addThingCorrectBody := &bytes.Buffer{}
	addThingCorrectWriter := multipart.NewWriter(addThingCorrectBody)
	addThingCorrectWriter.WriteField("thing_id", strconv.Itoa(thingID))
	addThingCorrectWriter.CreateFormFile("files", gofakeit.Word())
	addThingCorrectContentType := addThingCorrectWriter.FormDataContentType()
	addThingCorrectWriter.Close()

	// Incorrect request for adding place image
	addPlaceIncorrectBody := &bytes.Buffer{}
	addPlaceIncorrectWriter := multipart.NewWriter(addPlaceIncorrectBody)
	addPlaceIncorrectWriter.WriteField("place_id", gofakeit.Word())
	addPlaceIncorrectContentType := addPlaceIncorrectWriter.FormDataContentType()
	addPlaceIncorrectWriter.Close()

	// Incorrect request for adding thing image
	addThingIncorrectBody := &bytes.Buffer{}
	addThingIncorrectWriter := multipart.NewWriter(addThingIncorrectBody)
	addThingIncorrectWriter.WriteField("thing_id", gofakeit.Word())
	addThingIncorrectContentType := addThingIncorrectWriter.FormDataContentType()
	addThingIncorrectWriter.Close()

	// Incorrect empty request
	emptyBody := &bytes.Buffer{}
	emptyWriter := multipart.NewWriter(emptyBody)
	emptyContentType := emptyWriter.FormDataContentType()
	emptyWriter.Close()

	tests := []struct {
		name               string
		req                req
		resCode            int
		resBody            interface{}
		fileRepoMock       func(mc *minimock.Controller) FileRepository
		placeImageRepoMock func(mc *minimock.Controller) PlaceImageRepository
		thingImageRepoMock func(mc *minimock.Controller) ThingImageRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/images",
			},
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name: "negative case - incorrect request",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addPlaceIncorrectBody.Bytes(),
				contentType: addPlaceIncorrectContentType,
			},
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name: "negative case - incorrect request",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addThingIncorrectBody.Bytes(),
				contentType: addThingIncorrectContentType,
			},
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name: "negative case - empty request",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        emptyBody.Bytes(),
				contentType: emptyContentType,
			},
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name: "negative case - repository error (begin tx)",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addPlaceCorrectBody.Bytes(),
				contentType: addPlaceCorrectContentType,
			},
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "negative case - repository error (commit tx)",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addPlaceCorrectBody.Bytes(),
				contentType: addPlaceCorrectContentType,
			},
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceImageRequest, tx *sql.Tx) {
					assert.Equal(mc, placeID, req.PlaceID)
				}).Return(nil)

				mock.CommitTxMock.Return(testError)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "negative case - repository error (add image)",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addPlaceCorrectBody.Bytes(),
				contentType: addPlaceCorrectContentType,
			},
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceImageRequest, tx *sql.Tx) {
					assert.Equal(mc, placeID, req.PlaceID)
				}).Return(testError)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "positive case - add place image",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addPlaceCorrectBody.Bytes(),
				contentType: addPlaceCorrectContentType,
			},
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceImageRequest, tx *sql.Tx) {
					assert.Equal(mc, placeID, req.PlaceID)
				}).Return(nil)

				mock.CommitTxMock.Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "negative case - repository error (begin tx)",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addThingCorrectBody.Bytes(),
				contentType: addThingCorrectContentType,
			},
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "negative case - repository error (commit tx)",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addThingCorrectBody.Bytes(),
				contentType: addThingCorrectContentType,
			},
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingImageRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(nil)

				mock.CommitTxMock.Return(testError)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "negative case - repository error (add image)",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addThingCorrectBody.Bytes(),
				contentType: addThingCorrectContentType,
			},
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingImageRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(testError)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "positive case - add thing image",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addThingCorrectBody.Bytes(),
				contentType: addThingCorrectContentType,
			},
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingImageRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(nil)

				mock.CommitTxMock.Return(nil)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
		},
		{
			name: "negative case - save file error",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/images",
				body:        addThingCorrectBody.Bytes(),
				contentType: addThingCorrectContentType,
			},
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.SaveMock.Return(testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()

			fiberApp.Post("/v1/images", AddImageHandler(
				tt.fileRepoMock(mc),
				tt.thingImageRepoMock(mc),
				tt.placeImageRepoMock(mc),
			))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, bytes.NewReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
