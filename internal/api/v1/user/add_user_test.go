package user

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/user/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
)

func TestAddUserHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.AddUserRequest
	}

	var (
		id        = gofakeit.Number(1, 1000)
		username  = gofakeit.Username()
		password  = gofakeit.Word()
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodPost,
			route:  "/v1/users",
			body: &dto.AddUserRequest{
				Username: username,
				Password: password,
			},
			contentType: fiber.MIMEApplicationJSON,
		}
	)

	tests := []struct {
		name            string
		req             req
		resCode         int
		resBody         interface{}
		userRepoMock    func(mc *minimock.Controller) UserRepository
		authServiceMock func(mc *minimock.Controller) AuthService
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, reqUsername string, reqPassword string) {
					assert.Equal(mc, username, reqUsername)
					assert.NotEmpty(mc, reqPassword)
				}).Return(id, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GeneratePasswordHashMock.Return(password, nil)
				return mock
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/users",
			},
			resCode: fiber.StatusBadRequest,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/users",
				body: &dto.AddUserRequest{
					Password: password,
				},
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, reqUsername string, reqPassword string) {
					assert.Equal(mc, username, reqUsername)
					assert.NotEmpty(mc, reqPassword)
				}).Return(0, testError)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GeneratePasswordHashMock.Return(password, nil)
				return mock
			},
		},
		{
			name:    "negative case - auth service error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GeneratePasswordHashMock.Return("", testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Post("/v1/users", AddUserHandler(tt.authServiceMock(mc), tt.userRepoMock(mc)))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, helpers.ConvertDataToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
