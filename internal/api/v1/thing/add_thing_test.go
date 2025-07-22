package thing

import (
	"context"
	"database/sql"
	"net/http/httptest"
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

func TestAddThingHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.AddThingRequest
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
			method: fiber.MethodPost,
			route:  "/v1/things",
			body: &dto.AddThingRequest{
				PlaceID:     placeID,
				Title:       title,
				Description: description,
			},
			contentType: fiber.MIMEApplicationJSON,
		}

		repoRes = models.Thing{
			ID:          thingID,
			PlaceID:     placeID,
			Title:       title,
			Description: description,
			CreatedAt:   gofakeit.Date(),
			UpdatedAt:   gofakeit.Date(),
		}

		expectedRes = dto.ThingResponse{
			ID:          thingID,
			PlaceID:     placeID,
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

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(thingID, nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceThingRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, placeID, req.PlaceID)
				}).Return(nil)

				return mock
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/things",
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
				method:      fiber.MethodPost,
				route:       "/v1/things",
				contentType: fiber.MIMEApplicationJSON,
				body: &dto.AddThingRequest{
					Title:       title,
					Description: description,
				},
			},
			resCode: fiber.StatusBadRequest,
			resBody: []*dto.ValidateErrorResponse{
				{
					Field: "AddThingRequest.PlaceID",
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
				method:      fiber.MethodPost,
				route:       "/v1/things",
				contentType: fiber.MIMEApplicationJSON,
				body: &dto.AddThingRequest{
					PlaceID:     placeID,
					Description: description,
				},
			},
			resCode: fiber.StatusBadRequest,
			resBody: []*dto.ValidateErrorResponse{
				{
					Field: "AddThingRequest.Title",
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
			name:    "negative case - repository error (add thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(0, testError)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				return mocks.NewPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (add place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(thingID, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceThingRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, placeID, req.PlaceID)
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

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(thingID, nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) PlaceThingRepository {
				mock := mocks.NewPlaceThingRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceThingRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, placeID, req.PlaceID)
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
			fiberApp.Post("/v1/things", AddThingHandler(
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
