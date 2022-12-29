package place

import (
	"context"
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

func Test_GetPlaceThingsHandler(t *testing.T) {
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.IThingRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		placeID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/places/" + strconv.Itoa(placeID) + "/things",
		}

		thingRepoRes = []models.Thing{
			{
				ID:          gofakeit.Number(1, 1000),
				PlaceID:     placeID,
				Title:       gofakeit.Phrase(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date().String(),
				UpdatedAt:   gofakeit.Date().String(),
			},
			{
				ID:          gofakeit.Number(1, 1000),
				PlaceID:     placeID,
				Title:       gofakeit.Phrase(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date().String(),
				UpdatedAt:   gofakeit.Date().String(),
			},
		}

		expectedRes = dto.ThingsResponse{
			Things: []dto.ThingResponse{
				{
					ID:          thingRepoRes[0].ID,
					PlaceID:     thingRepoRes[0].PlaceID,
					Title:       thingRepoRes[0].Title,
					Description: thingRepoRes[0].Description,
					CreatedAt:   thingRepoRes[0].CreatedAt,
					UpdatedAt:   thingRepoRes[0].UpdatedAt,
				},
				{
					ID:          thingRepoRes[1].ID,
					PlaceID:     thingRepoRes[1].PlaceID,
					Title:       thingRepoRes[1].Title,
					Description: thingRepoRes[1].Description,
					CreatedAt:   thingRepoRes[1].CreatedAt,
					UpdatedAt:   thingRepoRes[1].UpdatedAt,
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
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)

				mock.GetAllByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(thingRepoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetAllByPlaceIDMock.Return(nil, testError)
				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/places/" + gofakeit.Word() + "/things",
			},
			resCode: fiber.StatusBadRequest,
			resBody: nil,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				return repoMocks.NewIThingRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingRepoMock(mc))

			fiberApp.Get("/v1/places/:id/things", GetPlaceThingsHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
