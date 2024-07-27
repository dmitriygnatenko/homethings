package auth

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i AuthService,UserRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

type (
	AuthService interface {
		GetClaims(fctx *fiber.Ctx) jwt.MapClaims
		IsCorrectPassword(password string, hash string) bool
		GenerateToken(user models.User) (string, error)
	}

	UserRepository interface {
		Get(ctx context.Context, username string) (*models.User, error)
	}
)

// @Router 		/api/v1/auth/check [get]
// @Success     200 {object} dto.UserResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Check auth
// @Tags  		Auth
// @Accept      json
// @Produce     json
func CheckAuthHandler(
	authService AuthService,
	userRepository UserRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		claims := authService.GetClaims(fctx)

		user, err := userRepository.Get(ctx, claims["name"].(string))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusForbidden, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ToUserResponse(*user))
	}
}
