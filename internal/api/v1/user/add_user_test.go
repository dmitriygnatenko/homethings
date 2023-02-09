package user

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	authMocks "git.dmitriygnatenko.ru/dima/homethings/internal/services/auth/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_AddUserHandler(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) interfaces.Auth
	type userRepoMockFunc func(mc *minimock.Controller) interfaces.UserRepository

	type req struct {
		method      string
		route       string
		body        *dto.AddUserRequest
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
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
		userRepoMock    userRepoMockFunc
		authServiceMock authServiceMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, reqUsername string, reqPassword string) {
					assert.Equal(mc, username, reqUsername)
					assert.NotEmpty(mc, reqPassword)
				}).Return(id, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				mock := authMocks.NewAuthMock(mc)
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
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				return repoMocks.NewUserRepositoryMock(mc)
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				return authMocks.NewAuthMock(mc)
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
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				return repoMocks.NewUserRepositoryMock(mc)
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				return authMocks.NewAuthMock(mc)
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, reqUsername string, reqPassword string) {
					assert.Equal(mc, username, reqUsername)
					assert.NotEmpty(mc, reqPassword)
				}).Return(0, testError)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				mock := authMocks.NewAuthMock(mc)
				mock.GeneratePasswordHashMock.Return(password, nil)
				return mock
			},
		},
		{
			name:    "negative case - auth service error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				return repoMocks.NewUserRepositoryMock(mc)
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				mock := authMocks.NewAuthMock(mc)
				mock.GeneratePasswordHashMock.Return("", testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.userRepoMock(mc), tt.authServiceMock(mc))

			fiberApp.Post("/v1/users", AddUserHandler(serviceProvider))

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
