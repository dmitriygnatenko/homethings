package image

import (
	"context"
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

func Test_DeletePlaceImageHandler(t *testing.T) {
	type placeImageRepoMockFunc func(mc *minimock.Controller) interfaces.IPlaceImageRepository
	type fileRepoMockFunc func(mc *minimock.Controller) interfaces.IFileRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		imageID   = gofakeit.Number(1, 1000)
		imageURL  = gofakeit.URL()
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/images/place/" + strconv.Itoa(imageID),
		}

		repoRes = &models.Image{
			Image: imageURL,
		}
	)

	tests := []struct {
		name               string
		req                req
		resCode            int
		resBody            interface{}
		placeImageRepoMock placeImageRepoMockFunc
		fileRepoMock       fileRepoMockFunc
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/images/place/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.IPlaceImageRepository {
				return repoMocks.NewIPlaceImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - bad request (image not exists)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.IPlaceImageRepository {
				mock := repoMocks.NewIPlaceImageRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.IPlaceImageRepository {
				mock := repoMocks.NewIPlaceImageRepositoryMock(mc)
				mock.GetMock.Return(nil, testError)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.IFileRepository {
				return repoMocks.NewIFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (update)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.IPlaceImageRepository {
				mock := repoMocks.NewIPlaceImageRepositoryMock(mc)

				mock.GetMock.Return(nil, nil)
				mock.DeleteMock.Return(testError)

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
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.IPlaceImageRepository {
				mock := repoMocks.NewIPlaceImageRepositoryMock(mc)

				mock.GetMock.Return(repoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

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
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.IPlaceImageRepository {
				mock := repoMocks.NewIPlaceImageRepositoryMock(mc)

				mock.GetMock.Return(repoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

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
			serviceProvider := sp.InitMock(tt.placeImageRepoMock(mc), tt.fileRepoMock(mc))

			fiberApp.Delete("/v1/images/place/:id", DeletePlaceImageHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
