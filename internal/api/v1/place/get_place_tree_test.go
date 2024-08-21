package place

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/place/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/test"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func TestGetPlaceTreeHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		id1 = uint64(gofakeit.Number(1, 1000))
		id2 = uint64(gofakeit.Number(1, 1000))
		id3 = uint64(gofakeit.Number(1, 1000))

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/places/tree",
		}

		placeRepoRes = []models.Place{
			{
				ID:        id1,
				Title:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
			{
				ID:        id2,
				ParentID:  sql.NullInt64{Valid: true, Int64: int64(id1)},
				Title:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
			{
				ID:        id3,
				ParentID:  sql.NullInt64{Valid: true, Int64: int64(id2)},
				Title:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}

		expectedRes = dto.PlaceTreeResponse{
			Places: []dto.PlaceTree{
				{
					Place: dto.PlaceResponse{
						ID:        id1,
						ParentID:  nil,
						Title:     placeRepoRes[0].Title,
						CreatedAt: placeRepoRes[0].CreatedAt.Format(layout),
						UpdatedAt: placeRepoRes[0].UpdatedAt.Format(layout),
					},
					NestedPlaces: []dto.PlaceTree{
						{
							Place: dto.PlaceResponse{
								ID:        id2,
								ParentID:  &id1,
								Title:     placeRepoRes[1].Title,
								CreatedAt: placeRepoRes[1].CreatedAt.Format(layout),
								UpdatedAt: placeRepoRes[1].UpdatedAt.Format(layout),
							},
							NestedPlaces: []dto.PlaceTree{
								{
									Place: dto.PlaceResponse{
										ID:        id3,
										ParentID:  &id2,
										Title:     placeRepoRes[2].Title,
										CreatedAt: placeRepoRes[2].CreatedAt.Format(layout),
										UpdatedAt: placeRepoRes[2].UpdatedAt.Format(layout),
									},
								},
							},
						},
					},
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
				mock.GetAllMock.Return(placeRepoRes, nil)
				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)
				mock.GetAllMock.Return(nil, testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/places/tree", GetPlaceTreeHandler(tt.placeRepoMock(mc)))

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
