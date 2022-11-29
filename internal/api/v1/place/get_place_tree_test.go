package place

import (
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
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

func Test_GetPlaceTreeHandler(t *testing.T) {
	type placeRepoMockFunc func(mc *minimock.Controller) interfaces.IPlaceRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		testError = errors.New(gofakeit.Phrase())

		id1 = gofakeit.Number(1, 1000)
		id2 = gofakeit.Number(1, 1000)
		id3 = gofakeit.Number(1, 1000)

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/places/tree",
		}

		placeRepoRes = []models.Place{
			{
				ID:        id1,
				Title:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date().String(),
				UpdatedAt: gofakeit.Date().String(),
			},
			{
				ID:        id2,
				ParentID:  sql.NullInt64{Valid: true, Int64: int64(id1)},
				Title:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date().String(),
				UpdatedAt: gofakeit.Date().String(),
			},
			{
				ID:        id3,
				ParentID:  sql.NullInt64{Valid: true, Int64: int64(id2)},
				Title:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date().String(),
				UpdatedAt: gofakeit.Date().String(),
			},
		}

		expectedRes = dto.PlaceTreeResponse{
			Places: []dto.PlaceTree{
				{
					Place: dto.PlaceResponse{
						ID:        id1,
						ParentID:  nil,
						Title:     placeRepoRes[0].Title,
						CreatedAt: placeRepoRes[0].CreatedAt,
						UpdatedAt: placeRepoRes[0].UpdatedAt,
					},
					NestedPlaces: []dto.PlaceTree{
						{
							Place: dto.PlaceResponse{
								ID:        id2,
								ParentID:  &id1,
								Title:     placeRepoRes[1].Title,
								CreatedAt: placeRepoRes[1].CreatedAt,
								UpdatedAt: placeRepoRes[1].UpdatedAt,
							},
							NestedPlaces: []dto.PlaceTree{
								{
									Place: dto.PlaceResponse{
										ID:        id3,
										ParentID:  &id2,
										Title:     placeRepoRes[2].Title,
										CreatedAt: placeRepoRes[2].CreatedAt,
										UpdatedAt: placeRepoRes[2].UpdatedAt,
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
		placeRepoMock placeRepoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				mock := repoMocks.NewIPlaceRepositoryMock(mc)
				mock.GetAllMock.Return(placeRepoRes, nil)
				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				mock := repoMocks.NewIPlaceRepositoryMock(mc)
				mock.GetAllMock.Return(nil, testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.placeRepoMock(mc))

			fiberApp.Get("/v1/places/tree", GetPlaceTreeHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), v1.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
