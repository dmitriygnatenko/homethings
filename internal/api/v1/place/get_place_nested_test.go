package place

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

func Test_GetNestedPlacesHandler(t *testing.T) {
	type placeRepoMockFunc func(mc *minimock.Controller) interfaces.IPlaceRepository

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
			route:  "/v1/places/" + strconv.Itoa(placeID) + "/nested",
		}

		repoRes = []models.Place{
			{
				ID:        placeID,
				ParentID:  sql.NullInt64{Int64: int64(placeID), Valid: true},
				Title:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}

		expectedRes = dto.PlacesResponse{
			Places: []dto.PlaceResponse{
				{
					ID:        repoRes[0].ID,
					ParentID:  &placeID,
					Title:     repoRes[0].Title,
					CreatedAt: repoRes[0].CreatedAt.Format(layout),
					UpdatedAt: repoRes[0].UpdatedAt.Format(layout),
				},
			},
		}
	)

	tests := []struct {
		name          string
		req           req
		resCode       int
		resBody       interface{}
		placeRepoMock placeRepoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				mock := repoMocks.NewIPlaceRepositoryMock(mc)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(repoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				mock := repoMocks.NewIPlaceRepositoryMock(mc)
				mock.GetNestedPlacesMock.Return(nil, testError)
				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/places/" + gofakeit.Word() + "/nested",
			},
			resCode: fiber.StatusBadRequest,
			resBody: nil,
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				return repoMocks.NewIPlaceRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.placeRepoMock(mc))

			fiberApp.Get("/v1/places/:id/nested", GetNestedPlacesHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
