package image

import (
	"context"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/dmitriygnatenko/homethings/internal/api/v1/image/mocks"
	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings/internal/models"
)

func TestDeletePlaceImageHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		imageID   = uint64(gofakeit.Number(1, 1000))
		imageURL  = gofakeit.URL()
		testError = gofakeit.Error()

		txMockFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/images/place/" + strconv.FormatUint(imageID, 10),
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
		tmMock             func(mc *minimock.Controller) TransactionManager
		placeImageRepoMock func(mc *minimock.Controller) PlaceImageRepository
		fileRepoMock       func(mc *minimock.Controller) FileRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/images/place/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil, testError)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(testError)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - file delete error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(repoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.DeleteMock.Return(testError)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(repoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				mock := mocks.NewFileRepositoryMock(mc)
				mock.DeleteMock.Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()

			fiberApp.Delete("/v1/images/place/:imageId", DeletePlaceImageHandler(
				tt.tmMock(mc),
				tt.fileRepoMock(mc),
				tt.placeImageRepoMock(mc),
			))

			fiberRes, _ := fiberApp.Test(
				httptest.NewRequest(tt.req.method, tt.req.route, nil),
				test.TestTimeout,
			)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)

			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
