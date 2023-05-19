package thing

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
	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		placeID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/things/place/" + strconv.Itoa(placeID),
		}

		thingRepoRes = []models.Thing{
			{
				ID:          gofakeit.Number(1, 1000),
				PlaceID:     placeID,
				Title:       gofakeit.Phrase(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date(),
				UpdatedAt:   gofakeit.Date(),
			},
			{
				ID:          gofakeit.Number(1, 1000),
				PlaceID:     placeID,
				Title:       gofakeit.Phrase(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date(),
				UpdatedAt:   gofakeit.Date(),
			},
		}

		thingTagRepoRes = []models.ThingTag{
			{
				ThingID: thingRepoRes[0].ID,
				Tag: models.Tag{
					ID:        gofakeit.Number(1, 1000),
					Title:     gofakeit.Phrase(),
					Style:     gofakeit.Phrase(),
					CreatedAt: gofakeit.Date(),
					UpdatedAt: gofakeit.Date(),
				},
			},
			{
				ThingID: thingRepoRes[1].ID,
				Tag: models.Tag{
					ID:        gofakeit.Number(1, 1000),
					Title:     gofakeit.Phrase(),
					Style:     gofakeit.Phrase(),
					CreatedAt: gofakeit.Date(),
					UpdatedAt: gofakeit.Date(),
				},
			},
		}

		expectedRes = dto.ThingsExtResponse{
			Things: []dto.ThingExtResponse{
				{
					ThingResponse: dto.ThingResponse{
						ID:          thingRepoRes[0].ID,
						PlaceID:     thingRepoRes[0].PlaceID,
						Title:       thingRepoRes[0].Title,
						Description: thingRepoRes[0].Description,
						CreatedAt:   thingRepoRes[0].CreatedAt.Format(layout),
						UpdatedAt:   thingRepoRes[0].UpdatedAt.Format(layout),
					},
					Tags: []dto.TagResponse{
						{
							ID:        thingTagRepoRes[0].ID,
							Title:     thingTagRepoRes[0].Title,
							Style:     thingTagRepoRes[0].Style,
							CreatedAt: thingTagRepoRes[0].CreatedAt.Format(layout),
							UpdatedAt: thingTagRepoRes[0].UpdatedAt.Format(layout),
						},
					},
				},
				{
					ThingResponse: dto.ThingResponse{
						ID:          thingRepoRes[1].ID,
						PlaceID:     thingRepoRes[1].PlaceID,
						Title:       thingRepoRes[1].Title,
						Description: thingRepoRes[1].Description,
						CreatedAt:   thingRepoRes[1].CreatedAt.Format(layout),
						UpdatedAt:   thingRepoRes[1].UpdatedAt.Format(layout),
					},
					Tags: []dto.TagResponse{
						{
							ID:        thingTagRepoRes[1].ID,
							Title:     thingTagRepoRes[1].Title,
							Style:     thingTagRepoRes[1].Style,
							CreatedAt: thingTagRepoRes[1].CreatedAt.Format(layout),
							UpdatedAt: thingTagRepoRes[1].UpdatedAt.Format(layout),
						},
					},
				},
			},
		}
	)

	tests := []struct {
		name             string
		req              req
		resCode          int
		resBody          interface{}
		thingRepoMock    func(mc *minimock.Controller) interfaces.ThingRepository
		thingTagRepoMock func(mc *minimock.Controller) interfaces.ThingTagRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetAllByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(thingRepoRes, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(thingTagRepoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - thing repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetAllByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing tag repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetAllByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(thingRepoRes, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/things/place/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			resBody: nil,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingRepoMock(mc), tt.thingTagRepoMock(mc))

			fiberApp.Get("/v1/things/place/:placeId", GetPlaceThingsHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
