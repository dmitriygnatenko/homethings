package user

import (
	"errors"
	"net/http/httptest"
	"testing"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
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

func Test_UpdateUserHandler(t *testing.T) {
	type userRepoMockFunc func(mc *minimock.Controller) interfaces.IUserRepository

	type req struct {
		method      string
		route       string
		body        *dto.UpdateUserRequest
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
		id        = gofakeit.Number(1, 1000)
		username  = gofakeit.Username()
		password  = gofakeit.Word()
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/users",
			body: &dto.UpdateUserRequest{
				Username: &username,
				Password: &password,
			},
			contentType: fiber.MIMEApplicationJSON,
		}
	)

	_ = correctReq
	_ = testError
	_ = id

	tests := []struct {
		name         string
		req          req
		resCode      int
		resBody      interface{}
		userRepoMock userRepoMockFunc
	}{
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/users",
			},
			resCode: fiber.StatusBadRequest,
			userRepoMock: func(mc *minimock.Controller) interfaces.IUserRepository {
				return repoMocks.NewIUserRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.userRepoMock(mc))

			fiberApp.Put("/v1/users", UpdateUserHandler(serviceProvider))

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
