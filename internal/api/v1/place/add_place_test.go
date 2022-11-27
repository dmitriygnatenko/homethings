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

func Test_AddPlaceHandler(t *testing.T) {
	type placeRepoMockFunc func(mc *minimock.Controller) interfaces.IPlaceRepository

	type req struct {
		method      string
		route       string
		body        *dto.AddPlaceRequest
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
		placeID   = gofakeit.Number(1, 1000)
		parentID  = gofakeit.Number(1, 1000)
		title     = gofakeit.Phrase()
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodPost,
			route:  "/v1/places",
			body: &dto.AddPlaceRequest{
				Title:    title,
				ParentID: &parentID,
			},
			contentType: fiber.MIMEApplicationJSON,
		}

		repoRes = models.Place{
			ID:        placeID,
			Title:     title,
			ParentID:  sql.NullInt64{Int64: int64(parentID), Valid: true},
			CreatedAt: gofakeit.Date().String(),
			UpdatedAt: gofakeit.Date().String(),
		}

		expectedRes = dto.PlaceResponse{
			ID:        placeID,
			ParentID:  &parentID,
			Title:     repoRes.Title,
			CreatedAt: repoRes.CreatedAt,
			UpdatedAt: repoRes.UpdatedAt,
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
				mock.AddMock.Return(placeID, nil)
				mock.GetMock.Return(&repoRes, nil)
				return mock
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/places",
			},
			resCode: fiber.StatusBadRequest,
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				return repoMocks.NewIPlaceRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without title",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/places",
				contentType: fiber.MIMEApplicationJSON,
				body:        &dto.AddPlaceRequest{},
			},
			resCode: fiber.StatusBadRequest,
			resBody: []*dto.ValidateErrorResponse{
				{
					Field: "AddPlaceRequest.Title",
					Tag:   "required",
				},
			},
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				return repoMocks.NewIPlaceRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (add place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				mock := repoMocks.NewIPlaceRepositoryMock(mc)
				mock.AddMock.Return(0, testError)
				return mock
			},
		},
		{
			name:    "negative case - repository error (get place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			placeRepoMock: func(mc *minimock.Controller) interfaces.IPlaceRepository {
				mock := repoMocks.NewIPlaceRepositoryMock(mc)
				mock.AddMock.Return(placeID, nil)
				mock.GetMock.Return(nil, testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.placeRepoMock(mc))

			fiberApp.Post("/v1/places", AddPlaceHandler(serviceProvider))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, helpers.ConvertDTOToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, v1.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
