package auth

import (
	"database/sql"
	"time"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// @Router 		/api/v1/auth/login [post]
// @Param       data body dto.LoginRequest true "Request body"
// @Success     200 {object} dto.LoginResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Login user
// @Tags  		Auth
// @Accept      json
// @Produce     json
func LoginHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.LoginRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		user, err := sp.GetUserRepository().Get(ctx, req.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return factory.CreateBadRequestResponse(fctx, nil)
			}

			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return factory.CreateBadRequestResponse(fctx, nil)
		}

		claims := jwt.MapClaims{
			"name": user.Username,
			"exp":  time.Now().Add(time.Duration(sp.GetEnvService().GetJWTLifetime()) * time.Second).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		encodedToken, err := token.SignedString([]byte(sp.GetEnvService().GetJWTSecretKey()))
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(dto.LoginResponse{Token: encodedToken})
	}
}
