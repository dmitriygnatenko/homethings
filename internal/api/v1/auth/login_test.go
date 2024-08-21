package auth

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/auth/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/test"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func TestLoginHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.LoginRequest
	}

	var (
		id           = gofakeit.Uint64()
		username     = gofakeit.Username()
		password     = gofakeit.Word()
		passwordHash = gofakeit.Word()
		token        = gofakeit.Word()
		testError    = gofakeit.Error()

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
		userRepoMock    func(mc *minimock.Controller) UserRepository
		authServiceMock func(mc *minimock.Controller) AuthService
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.LoginResponse{Token: token},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.IsCorrectPasswordMock.Return(true)
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
				route:  "/v1/auth/login",
				body: &dto.LoginRequest{
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
			name:    "negative case - generate token error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.IsCorrectPasswordMock.Expect(password, passwordHash).Return(true)
				mock.GenerateTokenMock.Expect(user).Return("", testError)

				return mock
			},
		},
		{
			name:    "negative case - incorrect password",
			req:     correctReq,
			resCode: fiber.StatusForbidden,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.IsCorrectPasswordMock.Expect(password, passwordHash).Return(false)
				return mock
			},
		},
		{
			name:    "negative case - repository error (get user)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(nil, testError)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
		},
		{
			name:    "negative case - user not found",
			req:     correctReq,
			resCode: fiber.StatusForbidden,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Post("/v1/auth/login", LoginHandler(tt.authServiceMock(mc), tt.userRepoMock(mc)))

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
