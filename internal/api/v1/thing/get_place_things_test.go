package thing

import (
	"context"
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

func TestGetPlaceThingsHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		placeID   = uint64(gofakeit.Number(1, 1000))
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/things/place/" + strconv.FormatUint(placeID, 10),
		}

		thingRepoRes = []models.Thing{
			{
				ID:          uint64(gofakeit.Number(1, 1000)),
				PlaceID:     placeID,
				Title:       gofakeit.Phrase(),
				Description: gofakeit.Phrase(),
				CreatedAt:   gofakeit.Date(),
				UpdatedAt:   gofakeit.Date(),
			},
			{
				ID:          uint64(gofakeit.Number(1, 1000)),
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
					ID:        uint64(gofakeit.Number(1, 1000)),
					Title:     gofakeit.Phrase(),
					Style:     gofakeit.Phrase(),
					CreatedAt: gofakeit.Date(),
					UpdatedAt: gofakeit.Date(),
				},
			},
			{
				ThingID: thingRepoRes[1].ID,
				Tag: models.Tag{
					ID:        uint64(gofakeit.Number(1, 1000)),
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
		thingRepoMock    func(mc *minimock.Controller) ThingRepository
		thingTagRepoMock func(mc *minimock.Controller) ThingTagRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetAllByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(thingRepoRes, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(thingTagRepoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - thing repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetAllByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing tag repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetAllByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(thingRepoRes, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
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
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/things/place/:placeId", GetPlaceThingsHandler(
				tt.thingRepoMock(mc),
				tt.thingTagRepoMock(mc),
			))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), test.TestTimeout)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
