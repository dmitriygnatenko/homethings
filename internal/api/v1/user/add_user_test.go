package user

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/user/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/test"
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
		username  = gofakeit.Username()
		password  = gofakeit.Word()
		hash      = gofakeit.Word()
		testError = gofakeit.Error()

		correctReq = req{
			method: fiber.MethodPost,
			route:  "/v1/users/",
			body: &dto.AddUserRequest{
				Username: username,
				Password: password,
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

				return mock
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, u string, p string) {
					assert.Equal(mc, hash, p)
					assert.Equal(mc, username, u)
				}).Return(gofakeit.Uint64(), nil)

				return mock
			},
		},
		{
			name:    "negative case - user repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			authService: func(mc *minimock.Controller) AuthService {
				mock := mocks.NewAuthServiceMock(mc)

				mock.GeneratePasswordHashMock.Inspect(func(p string) {
					assert.Equal(mc, password, p)
				}).Return(hash, nil)

				return mock
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, u string, p string) {
					assert.Equal(mc, hash, p)
					assert.Equal(mc, username, u)
				}).Return(0, testError)

				return mock
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
			name: "negative case - bad request",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/users/",
				body:        &dto.AddUserRequest{},
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
				method: fiber.MethodPost,
				route:  "/v1/users/",
				body:   &dto.AddUserRequest{},
			},
			resCode: fiber.StatusBadRequest,
			authService: func(mc *minimock.Controller) AuthService {
				return mocks.NewAuthServiceMock(mc)
			},
			userRepoMock: func(mc *minimock.Controller) UserRepository {
				return mocks.NewUserRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Post("/v1/users", AddUserHandler(tt.authService(mc), tt.userRepoMock(mc)))

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
