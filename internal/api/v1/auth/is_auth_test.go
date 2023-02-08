package auth

import (
	"net/http/httptest"
	"testing"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_CheckAuthHandler(t *testing.T) {
	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		correctReq = req{
			method:      fiber.MethodGet,
			route:       "/v1/auth/check",
			contentType: fiber.MIMEApplicationJSON,
		}

		expectedRes = dto.EmptyResponse{}
	)

	tests := []struct {
		name    string
		req     req
		resCode int
		resBody interface{}
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock()

			fiberApp.Get("/v1/auth/check", CheckAuthHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
