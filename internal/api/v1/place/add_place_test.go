package place

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/dmitriygnatenko/homethings/internal/api/v1/place/mocks"
	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings/internal/models"
)

func TestAddPlaceHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.AddPlaceRequest
	}

	var (
		placeID   = uint64(gofakeit.Number(1, 1000))
		parentID  = uint64(gofakeit.Number(1, 1000))
		title     = gofakeit.Phrase()
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

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
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		expectedRes = dto.PlaceResponse{
			ID:        placeID,
			ParentID:  &parentID,
			Title:     repoRes.Title,
			CreatedAt: repoRes.CreatedAt.Format(layout),
			UpdatedAt: repoRes.UpdatedAt.Format(layout),
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

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, parentID, uint64(req.ParentID.Int64))
				}).Return(placeID, nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(&repoRes, nil)

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
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				return mocks.NewPlaceRepositoryMock(mc)
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
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				return mocks.NewPlaceRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (add place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, parentID, uint64(req.ParentID.Int64))
				}).Return(0, testError)

				return mock
			},
		},
		{
			name:    "negative case - repository error (get place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) PlaceRepository {
				mock := mocks.NewPlaceRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddPlaceRequest) {
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, parentID, uint64(req.ParentID.Int64))
				}).Return(placeID, nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, testError)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Post("/v1/places", AddPlaceHandler(tt.placeRepoMock(mc)))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, test.ConvertDataToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, test.TestTimeout)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)

			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
