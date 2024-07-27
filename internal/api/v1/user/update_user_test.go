package user

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/user/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/auth"
)

func TestUpdateUserHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.UpdateUserRequest
	}

	var (
		id          = gofakeit.Number(1, 1000)
		username    = gofakeit.Username()
		password    = gofakeit.Word()
		newUsername = gofakeit.Username()
		newPassword = gofakeit.Word()
		testError   = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/users",
			body: &dto.UpdateUserRequest{
				Username: &newUsername,
				Password: &newPassword,
			},
			contentType: fiber.MIMEApplicationJSON,
		}

		claims = jwt.MapClaims{
			auth.ClaimsKeyName: username,
		}

		user = models.User{
			ID:       id,
			Username: username,
			Password: password,
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
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateUserRequest) {
					assert.Equal(mc, id, req.ID)
					assert.Equal(mc, newUsername, req.Username.String)
					assert.Equal(mc, newPassword, req.Password.String)
				}).Return(nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Expect(newPassword).Return(newPassword, nil)
				mock.GetClaimsMock.Return(claims)

				return mock
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
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
				method:      fiber.MethodPut,
				route:       "/v1/users",
				body:        &dto.UpdateUserRequest{},
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
			name:    "negative case - auth service error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GeneratePasswordHashMock.Expect(newPassword).Return("", testError)
				return mock
			},
		},
		{
			name:    "negative case - repository error (update)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateUserRequest) {
					assert.Equal(mc, id, req.ID)
					assert.Equal(mc, newUsername, req.Username.String)
					assert.Equal(mc, newPassword, req.Password.String)
				}).Return(testError)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Expect(newPassword).Return(newPassword, nil)
				mock.GetClaimsMock.Return(claims)

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
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Expect(newPassword).Return(newPassword, nil)
				mock.GetClaimsMock.Return(claims)

				return mock
			},
		},
		{
			name:    "negative case - repository error (user not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Expect(newPassword).Return(newPassword, nil)
				mock.GetClaimsMock.Return(claims)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Put("/v1/users", UpdateUserHandler(tt.authServiceMock(mc), tt.userRepoMock(mc)))

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
