package auth

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

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/auth/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func TestCheckAuthHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		username  = gofakeit.Username()
		testError = gofakeit.Error()

		correctReq = req{
			method:      fiber.MethodGet,
			route:       "/v1/auth/check",
			contentType: fiber.MIMEApplicationJSON,
		}

		claims = jwt.MapClaims{
			"name": username,
		}

		user = models.User{
			Username: username,
		}

		expectedRes = dto.UserResponse{Username: username}
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
			resBody: expectedRes,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(_ context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GetClaimsMock.Return(claims)
				return mock
			},
		},
		{
			name:    "negative case - user not found",
			req:     correctReq,
			resCode: fiber.StatusForbidden,
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(_ context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)
				mock.GetClaimsMock.Return(claims)
				return mock
			},
		},
		{
			name:    "negative case - repository error",
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
			fiberApp.Get("/v1/auth/check", CheckAuthHandler(tt.authServiceMock(mc), tt.userRepoMock(mc)))

			fiberRes, _ := fiberApp.Test(
				httptest.NewRequest(tt.req.method, tt.req.route, nil),
				test.TestTimeout,
			)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)

			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
