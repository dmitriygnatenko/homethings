package place

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

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/place/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func TestDeletePlaceHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		placeID       = uint64(gofakeit.Number(1, 1000))
		thingID       = uint64(gofakeit.Number(1, 1000))
		placeImageID  = uint64(gofakeit.Number(1, 1000))
		thingImageID  = uint64(gofakeit.Number(1, 1000))
		placeImageURL = gofakeit.URL()
		thingImageURL = gofakeit.URL()
		testError     = gofakeit.Error()

		txMockFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/places/" + strconv.FormatUint(placeID, 10),
		}

		things = []models.Thing{
			{
				ID:          thingID,
				PlaceID:     placeID,
				Title:       gofakeit.Word(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date(),
				UpdatedAt:   gofakeit.Date(),
			},
		}

		placeImages = []models.Image{
			{
				ID:    placeImageID,
				Image: placeImageURL,
			},
		}

		thingImages = []models.Image{
			{
				ID:    thingImageID,
				Image: thingImageURL,
			},
		}
	)

	tests := []struct {
		name                      string
		req                       req
		resCode                   int
		resBody                   interface{}
		tmMock                    func(mc *minimock.Controller) TransactionManager
		placeRepoMock             func(mc *minimock.Controller) PlaceRepository
		thingRepoMock             func(mc *minimock.Controller) ThingRepository
		placeImageRepoMock        func(mc *minimock.Controller) PlaceImageRepository
		thingImageRepoMock        func(mc *minimock.Controller) ThingImageRepository
		placeThingRepoMock        func(mc *minimock.Controller) PlaceThingRepository
		thingTagRepoMock          func(mc *minimock.Controller) ThingTagRepository
		thingNotificationRepoMock func(mc *minimock.Controller) ThingNotificationRepository
		fileRepoMock              func(mc *minimock.Controller) FileRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/places/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				return mocks.NewPlaceRepositoryMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
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
			name:    "negative case - bad request (place not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
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
			name:    "negative case - repository error (get place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
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
			name:    "negative case - repository error (get nested places)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
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
			name:    "negative case - not empty nested places",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return([]models.Place{{ID: gofakeit.Uint64()}}, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
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
			name:    "negative case - repository error (get place images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

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
			name:    "negative case - repository error (get thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

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
			name:    "negative case - repository error (get thing images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

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
			name:    "negative case - repository error (delete place images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(testError)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

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
			name:    "negative case - repository error (delete thing images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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
			name:    "negative case - repository error (delete place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(testError)

				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
				}).Return(nil)

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
			name:    "negative case - repository error (delete thing tag)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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
			name:    "negative case - repository error (delete thing notifications)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

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
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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
			name:    "negative case - repository error (delete place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(testError)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

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
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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
			name:    "negative case - repository error (delete place image files)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

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
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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

				mock.DeleteMock.Inspect(func(path string) {
					assert.Equal(mc, placeImageURL, path)
				}).Return(testError)

				return mock
			},
		},
		{
			name:    "negative case - repository error (delete thing image files)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

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
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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

				mock.DeleteMock.Set(func(path string) (err error) {
					if path == placeImageURL {
						return nil
					}

					return testError
				})

				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(things, nil)

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
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImages, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingImageID, id)
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

			fiberApp.Delete("/v1/places/:placeId", DeletePlaceHandler(
				tt.tmMock(mc),
				tt.placeRepoMock(mc),
				tt.thingRepoMock(mc),
				tt.placeImageRepoMock(mc),
				tt.thingImageRepoMock(mc),
				tt.placeThingRepoMock(mc),
				tt.thingTagRepoMock(mc),
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
