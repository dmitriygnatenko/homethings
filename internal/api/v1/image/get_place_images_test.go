package image

import (
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

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

func Test_GetPlaceImagesHandler(t *testing.T) {
	type placeImageRepoMockFunc func(mc *minimock.Controller) interfaces.PlaceImageRepository
	type thingImageRepoMockFunc func(mc *minimock.Controller) interfaces.ThingImageRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		placeID   = gofakeit.Number(1, 1000)
		thingID   = gofakeit.Number(1, 1000)
		date1     = gofakeit.Date()
		date2     = date1.Add(time.Hour)
		date3     = date2.Add(time.Hour)
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/images/place/" + strconv.Itoa(placeID),
		}

		placeImageRepoRes = []models.Image{
			{
				ID:        gofakeit.Number(1, 1000),
				Image:     gofakeit.URL(),
				PlaceID:   sql.NullInt64{Valid: true, Int64: int64(placeID)},
				CreatedAt: date1,
			},
			{
				ID:        gofakeit.Number(1, 1000),
				Image:     gofakeit.URL(),
				PlaceID:   sql.NullInt64{Valid: true, Int64: int64(placeID)},
				CreatedAt: date2,
			},
		}

		thingImageRepoRes = []models.Image{
			{
				ID:        gofakeit.Number(1, 1000),
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
		thingImageRepoMock thingImageRepoMockFunc
		placeImageRepoMock placeImageRepoMockFunc
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/images/place/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				return repoMocks.NewPlaceImageRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - place repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, testError)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, testError)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingImageRepoRes, nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingImageRepoMock(mc), tt.placeImageRepoMock(mc))

			fiberApp.Get("/v1/images/place/:placeId", GetPlaceImagesHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
