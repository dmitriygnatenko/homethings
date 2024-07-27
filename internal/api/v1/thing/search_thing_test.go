package thing

import (
	"context"
	"errors"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/thing/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func TestSearchThingHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		search    = gofakeit.LetterN(10)
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/things/search/" + url.QueryEscape(search),
		}

		thingRepoRes = []models.Thing{
			{
				ID:          gofakeit.Number(1, 1000),
				PlaceID:     gofakeit.Number(1, 1000),
				Title:       gofakeit.Phrase(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date(),
				UpdatedAt:   gofakeit.Date(),
			},
			{
				ID:          gofakeit.Number(1, 1000),
				PlaceID:     gofakeit.Number(1, 1000),
				Title:       gofakeit.Phrase(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date(),
				UpdatedAt:   gofakeit.Date(),
			},
		}

		expectedRes = dto.ThingsResponse{
			Things: []dto.ThingResponse{
				{
					ID:          thingRepoRes[0].ID,
					PlaceID:     thingRepoRes[0].PlaceID,
					Title:       thingRepoRes[0].Title,
					Description: thingRepoRes[0].Description,
					CreatedAt:   thingRepoRes[0].CreatedAt.Format(layout),
					UpdatedAt:   thingRepoRes[0].UpdatedAt.Format(layout),
				},
				{
					ID:          thingRepoRes[1].ID,
					PlaceID:     thingRepoRes[1].PlaceID,
					Title:       thingRepoRes[1].Title,
					Description: thingRepoRes[1].Description,
					CreatedAt:   thingRepoRes[1].CreatedAt.Format(layout),
					UpdatedAt:   thingRepoRes[1].UpdatedAt.Format(layout),
				},
			},
		}
	)

	tests := []struct {
		name          string
		req           req
		resCode       int
		resBody       interface{}
		thingRepoMock func(mc *minimock.Controller) ThingRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/things/search/" + url.QueryEscape(gofakeit.LetterN(2)),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/things/search/" + url.QueryEscape(gofakeit.LetterN(10)+":"),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.SearchMock.Inspect(func(ctx context.Context, s string) {
					assert.Equal(mc, search, s)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.SearchMock.Inspect(func(ctx context.Context, s string) {
					assert.Equal(mc, search, s)
				}).Return(thingRepoRes, nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/things/search/:search", SearchThingHandler(tt.thingRepoMock(mc)))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
