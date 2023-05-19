package image

import (
	"context"
	"database/sql"
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

func Test_GetThingImagesHandler(t *testing.T) {
	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		thingID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/images/thing/" + strconv.Itoa(thingID),
		}

		imageRepoRes = []models.Image{
			{
				ID:        gofakeit.Number(1, 1000),
				Image:     gofakeit.URL(),
				CreatedAt: gofakeit.Date(),
				ThingID:   sql.NullInt64{Valid: true, Int64: int64(thingID)},
			},
			{
				ID:        gofakeit.Number(1, 1000),
				Image:     gofakeit.URL(),
				CreatedAt: gofakeit.Date(),
				ThingID:   sql.NullInt64{Valid: true, Int64: int64(thingID)},
			},
			{
				ID:        gofakeit.Number(1, 1000),
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
		thingImageRepoMock func(mc *minimock.Controller) interfaces.ThingImageRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/images/thing/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id int) {
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
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(imageRepoRes, nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingImageRepoMock(mc))

			fiberApp.Get("/v1/images/thing/:thingId", GetThingImagesHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
