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
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func Test_CheckAuthHandler(t *testing.T) {
	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
		username  = gofakeit.Username()
		testError = errors.New(gofakeit.Phrase())

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
		userRepoMock    func(mc *minimock.Controller) interfaces.UserRepository
		authServiceMock func(mc *minimock.Controller) interfaces.Auth
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			userRepoMock: func(mc *minimock.Controller) interfaces.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, reqUsername string) {
					assert.Equal(mc, username, reqUsername)
				}).Return(&user, nil)

				return mock
			},
			authServiceMock: func(mc *minimock.Controller) interfaces.Auth {
				mock := authMocks.NewAuthMock(mc)
				mock.GetClaimsMock.Return(claims)
				return mock
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
				mock := authMocks.NewAuthMock(mc)
				mock.GetClaimsMock.Return(claims)
				return mock
			},
		},
		{
			name:    "negative case - repository error",
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
				mock := authMocks.NewAuthMock(mc)
				mock.GetClaimsMock.Return(claims)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.userRepoMock(mc), tt.authServiceMock(mc))

			fiberApp.Get("/v1/auth/check", CheckAuthHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
