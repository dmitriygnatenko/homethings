package image

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/image/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func TestGetPlaceImagesHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		placeID   = uint64(gofakeit.Number(1, 1000))
		thingID   = uint64(gofakeit.Number(1, 1000))
		date1     = gofakeit.Date()
		date2     = date1.Add(time.Hour)
		date3     = date2.Add(time.Hour)
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/images/place/" + strconv.FormatUint(placeID, 10),
		}

		placeImageRepoRes = []models.Image{
			{
				ID:        gofakeit.Uint64(),
				Image:     gofakeit.URL(),
				PlaceID:   sql.NullInt64{Valid: true, Int64: int64(placeID)},
				CreatedAt: date1,
			},
			{
				ID:        gofakeit.Uint64(),
				Image:     gofakeit.URL(),
				PlaceID:   sql.NullInt64{Valid: true, Int64: int64(placeID)},
				CreatedAt: date2,
			},
		}

		thingImageRepoRes = []models.Image{
			{
				ID:        gofakeit.Uint64(),
				Image:     gofakeit.URL(),
				ThingID:   sql.NullInt64{Valid: true, Int64: int64(thingID)},
				CreatedAt: date3,
			},
		}

		expectedRes = dto.ImagesResponse{
			Images: []dto.ImageResponse{
				{
					ID:        thingImageRepoRes[0].ID,
					Image:     thingImageRepoRes[0].Image,
					CreatedAt: thingImageRepoRes[0].CreatedAt.Format(layout),
					ThingID:   &thingID,
				},
				{
					ID:        placeImageRepoRes[1].ID,
					Image:     placeImageRepoRes[1].Image,
					CreatedAt: placeImageRepoRes[1].CreatedAt.Format(layout),
					PlaceID:   &placeID,
				},
				{
					ID:        placeImageRepoRes[0].ID,
					Image:     placeImageRepoRes[0].Image,
					CreatedAt: placeImageRepoRes[0].CreatedAt.Format(layout),
					PlaceID:   &placeID,
				},
			},
		}
	)

	tests := []struct {
		name               string
		req                req
		resCode            int
		resBody            interface{}
		thingImageRepoMock func(mc *minimock.Controller) ThingImageRepository
		placeImageRepoMock func(mc *minimock.Controller) PlaceImageRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/images/place/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				return mocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - place repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
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
		},
		{
			name:    "negative case - thing repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			placeImageRepoMock: func(mc *minimock.Controller) PlaceImageRepository {
				mock := mocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImageRepoRes, nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(thingImageRepoRes, nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()

			fiberApp.Get("/v1/images/place/:placeId", GetPlaceImagesHandler(
				tt.thingImageRepoMock(mc),
				tt.placeImageRepoMock(mc),
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
