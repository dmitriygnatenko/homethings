package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	authMocks "git.dmitriygnatenko.ru/dima/homethings/internal/services/auth/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_LoginHandler(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) interfaces.Auth
	type userRepoMockFunc func(mc *minimock.Controller) interfaces.UserRepository

	type req struct {
		method      string
		route       string
		body        *dto.LoginRequest
		contentType string
	}

	var (
		mc           = minimock.NewController(t)
		id           = gofakeit.Number(1, 1000)
		username     = gofakeit.Username()
		password     = gofakeit.Word()
		passwordHash = gofakeit.Word()
		token        = gofakeit.Word()
		testError    = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodPost,
			route:  "/v1/auth/login",
			body: &dto.LoginRequest{
				Username: username,
				Password: password,
			},
			contentType: fiber.MIMEApplicationJSON,
		}

		user = models.User{
			ID:       id,
			Username: username,
			Password: passwordHash,
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
			resBody: dto.LoginResponse{Token: token},
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				mock := authMocks.NewAuthMock(mc)

				mock.IsCorrectPasswordMock.Expect(password, passwordHash).Return(true)
				mock.GenerateTokenMock.Expect(user).Return(token, nil)

				return mock
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/auth/login",
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
				route:  "/v1/auth/login",
				body: &dto.LoginRequest{
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
			name:    "negative case - generate token error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				mock := authMocks.NewAuthMock(mc)

				mock.IsCorrectPasswordMock.Expect(password, passwordHash).Return(true)
				mock.GenerateTokenMock.Expect(user).Return("", testError)

				return mock
			},
		},
		{
			name:    "negative case - incorrect password",
			req:     correctReq,
			resCode: fiber.StatusForbidden,
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				mock := authMocks.NewAuthMock(mc)
				mock.IsCorrectPasswordMock.Expect(password, passwordHash).Return(false)
				return mock
			},
		},
		{
			name:    "negative case - repository error (get user)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(nil, testError)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				return authMocks.NewAuthMock(mc)
			},
		},
		{
			name:    "negative case - user not found",
			req:     correctReq,
			resCode: fiber.StatusForbidden,
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				return authMocks.NewAuthMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.userRepoMock(mc), tt.authServiceMock(mc))

			fiberApp.Post("/v1/auth/login", LoginHandler(serviceProvider))

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
