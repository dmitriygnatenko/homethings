package thing

import (
	"errors"
	"net/http/httptest"
	"net/url"
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

func Test_SearchThingHandler(t *testing.T) {
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.IThingRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc = minimock.NewController(t)

		search           = gofakeit.LetterN(10)
		incorrectSearch  = gofakeit.LetterN(2)
		incorrectSearch2 = gofakeit.LetterN(10) + ":"
		testError        = errors.New(gofakeit.Phrase())
		layout           = "2006-01-02 15:04:05"

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
		thingRepoMock thingRepoMockFunc
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/things/search/" + url.QueryEscape(incorrectSearch),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				return repoMocks.NewIThingRepositoryMock(mc)
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/things/search/" + url.QueryEscape(incorrectSearch2),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				return repoMocks.NewIThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.SearchMock.Return(nil, testError)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.SearchMock.Return(thingRepoRes, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingRepoMock(mc))

			fiberApp.Get("/v1/things/search/:search", SearchThingHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
