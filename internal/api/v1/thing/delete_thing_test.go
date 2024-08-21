package thing

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/thing/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/test"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func TestDeleteThingHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		thingID   = uint64(gofakeit.Number(1, 1000))
		imageID   = uint64(gofakeit.Number(1, 1000))
		imageURL  = gofakeit.URL()
		testError = gofakeit.Error()

		txMockFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/things/" + strconv.FormatUint(thingID, 10),
		}

		repoImagesRes = []models.Image{
			{
				ID:    imageID,
				Image: imageURL,
			},
		}
	)

	tests := []struct {
		name                      string
		req                       req
		resCode                   int
		resBody                   interface{}
		tmMock                    func(mc *minimock.Controller) TransactionManager
		thingRepoMock             func(mc *minimock.Controller) ThingRepository
		placeThingRepoMock        func(mc *minimock.Controller) PlaceThingRepository
		thingImageRepoMock        func(mc *minimock.Controller) ThingImageRepository
		thingTagRepoMock          func(mc *minimock.Controller) ThingTagRepository
		thingNotificationRepoMock func(mc *minimock.Controller) ThingNotificationRepository
		fileRepoMock              func(mc *minimock.Controller) FileRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/things/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - bad request (thing not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(testError)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(repoImagesRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete thing tags)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(repoImagesRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(testError)

				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete notification)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(repoImagesRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(testError)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) FileRepository {
				return mocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(testError)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(repoImagesRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

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
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(repoImagesRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
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
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(repoImagesRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, imageID, id)
				}).Return(nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
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

			fiberApp.Delete("/v1/things/:thingId", DeleteThingHandler(
				tt.tmMock(mc),
				tt.thingRepoMock(mc),
				tt.thingTagRepoMock(mc),
				tt.placeThingRepoMock(mc),
				tt.thingImageRepoMock(mc),
				tt.thingNotificationRepoMock(mc),
				tt.fileRepoMock(mc),
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
