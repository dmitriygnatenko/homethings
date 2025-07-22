package image

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

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/image/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func TestGetThingImagesHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		thingID   = uint64(gofakeit.Number(1, 1000))
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/images/thing/" + strconv.FormatUint(thingID, 10),
		}

		imageRepoRes = []models.Image{
			{
				ID:        gofakeit.Uint64(),
				Image:     gofakeit.URL(),
				CreatedAt: gofakeit.Date(),
				ThingID:   sql.NullInt64{Valid: true, Int64: int64(thingID)},
			},
			{
				ID:        gofakeit.Uint64(),
				Image:     gofakeit.URL(),
				CreatedAt: gofakeit.Date(),
				ThingID:   sql.NullInt64{Valid: true, Int64: int64(thingID)},
			},
			{
				ID:        gofakeit.Uint64(),
				Image:     gofakeit.URL(),
				CreatedAt: gofakeit.Date(),
				ThingID:   sql.NullInt64{Valid: true, Int64: int64(thingID)},
			},
		}

		expectedRes = dto.ImagesResponse{
			Images: []dto.ImageResponse{
				{
					ID:        imageRepoRes[0].ID,
					Image:     imageRepoRes[0].Image,
					CreatedAt: imageRepoRes[0].CreatedAt.Format(layout),
					ThingID:   &thingID,
				},
				{
					ID:        imageRepoRes[1].ID,
					Image:     imageRepoRes[1].Image,
					CreatedAt: imageRepoRes[1].CreatedAt.Format(layout),
					ThingID:   &thingID,
				},
				{
					ID:        imageRepoRes[2].ID,
					Image:     imageRepoRes[2].Image,
					CreatedAt: imageRepoRes[2].CreatedAt.Format(layout),
					ThingID:   &thingID,
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
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/images/thing/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				return mocks.NewThingImageRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingImageRepoMock: func(mc *minimock.Controller) ThingImageRepository {
				mock := mocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(imageRepoRes, nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()

			fiberApp.Get("/v1/images/thing/:thingId", GetThingImagesHandler(tt.thingImageRepoMock(mc)))

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
