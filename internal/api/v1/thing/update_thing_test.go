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

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/thing/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func TestUpdateThingHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.UpdateThingRequest
	}

	var (
		placeID     = uint64(gofakeit.Number(1, 1000))
		thingID     = uint64(gofakeit.Number(1, 1000))
		title       = gofakeit.Phrase()
		description = gofakeit.Phrase()
		testError   = gofakeit.Error()
		layout      = "2006-01-02 15:04:05"

		txMockFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/things/" + strconv.FormatUint(thingID, 10),
			body: &dto.UpdateThingRequest{
				PlaceID:     placeID,
				Title:       title,
				Description: description,
			},
			contentType: fiber.MIMEApplicationJSON,
		}

		repoResBeforeUpdate = models.Thing{
			ID:          thingID,
			Title:       gofakeit.Phrase(),
			Description: gofakeit.Phrase(),
			CreatedAt:   gofakeit.Date(),
			UpdatedAt:   gofakeit.Date(),
		}

		placeThingRepoResBeforeUpdate = models.PlaceThing{
			PlaceID:   uint64(gofakeit.Number(1, 1000)),
			ThingID:   thingID,
			CreatedAt: gofakeit.Date(),
		}

		repoRes = models.Thing{
			ID:          thingID,
			Title:       title,
			Description: description,
			CreatedAt:   gofakeit.Date(),
			UpdatedAt:   gofakeit.Date(),
		}

		expectedRes = dto.ThingResponse{
			ID:          thingID,
			Title:       title,
			Description: description,
			CreatedAt:   repoRes.CreatedAt.Format(layout),
			UpdatedAt:   repoRes.UpdatedAt.Format(layout),
		}
	)

	tests := []struct {
		name               string
		req                req
		resCode            int
		resBody            interface{}
		tmMock             func(mc *minimock.Controller) TransactionManager
		thingRepoMock      func(mc *minimock.Controller) ThingRepository
		placeThingRepoMock func(mc *minimock.Controller) PlaceThingRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(nil)

				mock.GetMock.Set(func(ctx context.Context, id uint64) (*models.Thing, error) {
					assert.Equal(mc, thingID, id)
					if mock.GetAfterCounter() == 0 {
						return &repoResBeforeUpdate, nil
					}
					return &repoRes, nil
				})

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(&placeThingRepoResBeforeUpdate, nil)

				mock.UpdatePlaceMock.Inspect(func(ctx context.Context, req models.UpdatePlaceThingRequest) {
					assert.Equal(mc, placeID, req.PlaceID)
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(nil)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodPut,
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
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/" + strconv.FormatUint(thingID, 10),
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
		},
		{
			name: "negative case - request without place_id",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/things/" + strconv.FormatUint(thingID, 10),
				contentType: fiber.MIMEApplicationJSON,
				body: &dto.UpdateThingRequest{
					Title: title,
				},
			},
			resCode: fiber.StatusBadRequest,
			resBody: []*dto.ValidateErrorResponse{
				{
					Field: "UpdateThingRequest.PlaceID",
					Tag:   "required",
				},
			},
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without title",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/things/" + strconv.FormatUint(thingID, 10),
				contentType: fiber.MIMEApplicationJSON,
				body: &dto.UpdateThingRequest{
					PlaceID: placeID,
				},
			},
			resCode: fiber.StatusBadRequest,
			resBody: []*dto.ValidateErrorResponse{
				{
					Field: "UpdateThingRequest.Title",
					Tag:   "required",
				},
			},
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
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
		},
		{
			name:    "negative case - repository error (get place thing)",
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
				}).Return(&repoResBeforeUpdate, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
		},
		{
			name:    "negative case - repository error (update thing)",
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
				}).Return(&repoResBeforeUpdate, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(testError)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(&placeThingRepoResBeforeUpdate, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error (update place)",
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
				}).Return(&repoResBeforeUpdate, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(&placeThingRepoResBeforeUpdate, nil)

				mock.UpdatePlaceMock.Inspect(func(ctx context.Context, req models.UpdatePlaceThingRequest) {
					assert.Equal(mc, placeID, req.PlaceID)
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(testError)

				return mock
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
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Set(func(ctx context.Context, thingID uint64) (*models.Thing, error) {
					if mock.GetAfterCounter() == 0 {
						return &repoResBeforeUpdate, nil
					}
					return nil, sql.ErrNoRows
				})

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(&placeThingRepoResBeforeUpdate, nil)

				mock.UpdatePlaceMock.Inspect(func(ctx context.Context, req models.UpdatePlaceThingRequest) {
					assert.Equal(mc, placeID, req.PlaceID)
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Put("/v1/things/:thingId", UpdateThingHandler(
				tt.tmMock(mc),
				tt.thingRepoMock(mc),
				tt.placeThingRepoMock(mc),
			))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, test.ConvertDataToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, test.TestTimeout)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)

			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
