package place

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

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/place/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func TestGetNestedPlacesHandler(t *testing.T) {
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
			route:  "/v1/places/" + strconv.FormatUint(placeID, 10) + "/nested",
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
		placeRepoMock func(mc *minimock.Controller) PlaceRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(repoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

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
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				return mocks.NewPlaceRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/places/:parentPlaceId/nested", GetNestedPlacesHandler(tt.placeRepoMock(mc)))

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
