package user

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/user/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
	"github.com/dmitriygnatenko/homethings-v1/internal/services/auth"
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
		username  = gofakeit.Username()
		password  = gofakeit.Word()
		hash      = gofakeit.Word()
		testError = gofakeit.Error()

		claims = jwt.MapClaims{
			auth.ClaimsKeyName: username,
		}

		user = models.User{
			ID:        gofakeit.Uint64(),
			Username:  username,
			Password:  password,
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/users/",
			body: &dto.UpdateUserRequest{
				Username: &username,
				Password: &password,
			},
			contentType: fiber.MIMEApplicationJSON,
		}
	)

	tests := []struct {
		name         string
		req          req
		resCode      int
		resBody      interface{}
		authService  func(mc *minimock.Controller) AuthService
		userRepoMock func(mc *minimock.Controller) UserRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			authService: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Inspect(func(p string) {
					assert.Equal(mc, password, p)
				}).Return(hash, nil)

				mock.GetClaimsMock.Return(claims)

				return mock
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, u string) {
					assert.Equal(mc, username, u)
				}).Return(&user, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateUserRequest) {
					assert.Equal(mc, hash, req.Password.String)
					assert.Equal(mc, username, req.Username.String)
				}).Return(nil)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/users/",
				body:        &dto.UpdateUserRequest{},
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			authService: func(mc *minimock.Controller) AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/users/",
				body:   &dto.UpdateUserRequest{},
			},
			resCode: fiber.StatusBadRequest,
			authService: func(mc *minimock.Controller) AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - auth service error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			authService: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Inspect(func(p string) {
					assert.Equal(mc, password, p)
				}).Return("", testError)

				return mock
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error - get user",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			authService: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Inspect(func(p string) {
					assert.Equal(mc, password, p)
				}).Return(hash, nil)

				mock.GetClaimsMock.Return(claims)

				return mock
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, u string) {
					assert.Equal(mc, username, u)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name:    "negative case - repository error - update user",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			authService: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Inspect(func(p string) {
					assert.Equal(mc, password, p)
				}).Return(hash, nil)

				mock.GetClaimsMock.Return(claims)

				return mock
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, u string) {
					assert.Equal(mc, username, u)
				}).Return(&user, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateUserRequest) {
					assert.Equal(mc, hash, req.Password.String)
					assert.Equal(mc, username, req.Username.String)
				}).Return(testError)

				return mock
			},
		},
		{
			name:    "negative case - user not found",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			authService: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Inspect(func(p string) {
					assert.Equal(mc, password, p)
				}).Return(hash, nil)

				mock.GetClaimsMock.Return(claims)

				return mock
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, u string) {
					assert.Equal(mc, username, u)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Put("/v1/users", UpdateUserHandler(tt.authService(mc), tt.userRepoMock(mc)))

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
