package thing

import (
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteThingHandler(t *testing.T) {
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.IThingRepository
	type placeThingRepoMockFunc func(mc *minimock.Controller) interfaces.IPlaceThingRepository
	type thingImageRepoMockFunc func(mc *minimock.Controller) interfaces.IThingImageRepository
	type fileRepoMockFunc func(mc *minimock.Controller) interfaces.IFileRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		thingID   = gofakeit.Number(1, 1000)
		imageURL  = gofakeit.URL()
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/things/" + strconv.Itoa(thingID),
		}

		repoImagesRes = []models.Image{
			{
				Image: imageURL,
			},
		}
	)

	tests := []struct {
		name               string
		req                req
		resCode            int
		resBody            interface{}
		thingRepoMock      thingRepoMockFunc
		placeThingRepoMock placeThingRepoMockFunc
		thingImageRepoMock thingImageRepoMockFunc
		fileRepoMock       fileRepoMockFunc
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/things/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				return repoMocks.NewIThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - bad request (thing not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				return repoMocks.NewIThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				return repoMocks.NewIThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (begin tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				return repoMocks.NewIThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(testError)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				return repoMocks.NewIThingImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				mock := repoMocks.NewIThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(nil, testError)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				mock := repoMocks.NewIThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				mock := repoMocks.NewIThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (commit tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				mock := repoMocks.NewIThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - file delete error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				mock := repoMocks.NewIThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				mock := repoMocks.NewIFileRepositoryMock(mc)
				mock.DeleteMock.Return(testError)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.IThingImageRepository {
				mock := repoMocks.NewIThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				mock := repoMocks.NewIFileRepositoryMock(mc)
				mock.DeleteMock.Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(
				tt.thingRepoMock(mc),
				tt.placeThingRepoMock(mc),
				tt.thingImageRepoMock(mc),
				tt.fileRepoMock(mc),
			)

			fiberApp.Delete("/v1/things/:id", DeleteThingHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
