package thing

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
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

func Test_AddThingHandler(t *testing.T) {
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.ThingRepository
	type placeThingRepoMockFunc func(mc *minimock.Controller) interfaces.PlaceThingRepository

	type req struct {
		method      string
		route       string
		body        *dto.AddThingRequest
		contentType string
	}

	var (
		mc          = minimock.NewController(t)
		placeID     = gofakeit.Number(1, 1000)
		thingID     = gofakeit.Number(1, 1000)
		title       = gofakeit.Phrase()
		description = gofakeit.Phrase()
		testError   = errors.New(gofakeit.Phrase())
		layout      = "2006-01-02 15:04:05"

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
		thingRepoMock      thingRepoMockFunc
		placeThingRepoMock placeThingRepoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingRequest, tx *sql.Tx) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, description, req.Description)
				}).Return(thingID, nil)

				mock.CommitTxMock.Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceThingRequest, tx *sql.Tx) {
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
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
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
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
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
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (begin tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (add thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.AddMock.Return(0, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (add place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.AddMock.Return(thingID, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.AddMock.Return(testError)
				return mock
			},
		},
		{
			name:    "negative case - repository error (commit tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.AddMock.Return(thingID, nil)
				mock.CommitTxMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.AddMock.Return(nil)
				return mock
			},
		},
		{
			name:    "negative case - repository error (get thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.AddMock.Return(thingID, nil)
				mock.CommitTxMock.Return(nil)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.AddMock.Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingRepoMock(mc), tt.placeThingRepoMock(mc))

			fiberApp.Post("/v1/things", AddThingHandler(serviceProvider))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, helpers.ConvertDataToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
