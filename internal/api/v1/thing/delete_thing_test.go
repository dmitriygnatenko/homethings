package thing

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteThingHandler(t *testing.T) {
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.IThingRepository
	type placeThingRepoMockFunc func(mc *minimock.Controller) interfaces.IPlaceThingRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		thingID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/things/" + strconv.Itoa(thingID),
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
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/things/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - bad request (thing not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (begin tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				return repoMocks.NewIPlaceThingRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(testError)
				return mock
			},
		},
		{
			name:    "negative case - repository error (delete thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
		},
		{
			name:    "negative case - commit tx error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			resBody: dto.ErrorResponse{Error: testError.Error()},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				return mock.CommitTxMock.Return(testError)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			thingRepoMock: func(mc *minimock.Controller) interfaces.IThingRepository {
				mock := repoMocks.NewIThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)
				mock.CommitTxMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.IPlaceThingRepository {
				mock := repoMocks.NewIPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.thingRepoMock(mc), tt.placeThingRepoMock(mc))

			fiberApp.Delete("/v1/things/:id", DeleteThingHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), v1.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
