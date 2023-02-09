package interfaces

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Auth interface {
	GeneratePasswordHash(password string) (string, error)
	IsCorrectPassword(password string, hash string) bool
	GenerateToken(user models.User) (string, error)
	GetClaims(fctx *fiber.Ctx) jwt.MapClaims
}
