package thing

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strconv"
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

func Test_UpdateThingHandler(t *testing.T) {
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.IThingRepository
	type placeThingRepoMockFunc func(mc *minimock.Controller) interfaces.IPlaceThingRepository

	type req struct {
		method string
		route  string
		body   *dto.UpdateThingRequest
	}

	var (
		mc          = minimock.NewController(t)
		placeID     = gofakeit.Number(1, 1000)
		thingID     = gofakeit.Number(1, 1000)
		title       = gofakeit.Phrase()
		description = gofakeit.Phrase()
		testError   = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/things/" + strconv.Itoa(thingID),
			body: &dto.UpdateThingRequest{
				PlaceID:     &placeID,
				Title:       &title,
				Description: &description,
			},
		}

		repoRes = models.Thing{
			ID:          thingID,
			Title:       title,
			Description: description,
			CreatedAt:   gofakeit.Date().String(),
			UpdatedAt:   gofakeit.Date().String(),
		}

		expectedRes = dto.ThingResponse{
			ID:          thingID,
			Title:       title,
			Description: description,
			CreatedAt:   repoRes.CreatedAt,
			UpdatedAt:   repoRes.UpdatedAt,
		}
	)

	tests := []struct {
		name               string
		req                req
		resCode            int
		resBody            interface{}
		thingRepoMock      thingRepoMockFunc
		placeThingRepoMock placeThingRepoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)

				mock.BeginTxMock.Return(nil, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingRequest, tx *sql.Tx) {
					assert.Equal(mc, title, req.Title.String)
					assert.Equal(mc, description, req.Description.String)
				}).Return(nil)

				mock.CommitTxMock.Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)

				mock.UpdatePlaceMock.Inspect(func(ctx context.Context, req models.UpdatePlaceThingRequest, tx *sql.Tx) {
					assert.Equal(mc, placeID, req.PlaceID)
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(nil)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				return repoMocks.NewIThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/" + strconv.Itoa(thingID),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				return repoMocks.NewIThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (begin tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (update thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.UpdateMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (update place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.UpdateMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.UpdatePlaceMock.Return(testError)
				return mock
			},
		},
		{
			name:    "negative case - repository error (commit tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.UpdateMock.Return(nil)
				mock.CommitTxMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.UpdatePlaceMock.Return(nil)
				return mock
			},
		},
		{
			name:    "negative case - repository error (get thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.BeginTxMock.Return(nil, nil)
				mock.UpdateMock.Return(nil)
				mock.CommitTxMock.Return(nil)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.UpdatePlaceMock.Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingRepoMock(mc), tt.placeThingRepoMock(mc))

			fiberApp.Put("/v1/things/:id", UpdateThingHandler(serviceProvider))

			var reqBody io.Reader
			if tt.req.body != nil {
				b, _ := json.Marshal(tt.req.body)
				reqBody = bytes.NewReader(b)
			}

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, reqBody)
			fiberReq.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			fiberRes, _ := fiberApp.Test(fiberReq, v1.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
