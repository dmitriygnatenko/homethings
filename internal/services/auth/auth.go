package auth

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.Auth -o ./mocks/ -s "_minimock.go"

import (
	"time"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCost    = bcrypt.DefaultCost
	defaultUserKey = "user"
)

type auth struct {
	env interfaces.Env
}

func Init(env interfaces.Env) (interfaces.Auth, error) {
	return auth{env: env}, nil
}

func (a auth) GeneratePasswordHash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (a auth) IsCorrectPassword(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (a auth) GetClaims(fctx *fiber.Ctx) jwt.MapClaims {
	jwtUser := fctx.Locals(defaultUserKey).(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)

	return claims
}

func (a auth) GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"name": user.Username,
		"exp":  time.Now().Add(time.Duration(a.env.GetJWTLifetime()) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.env.GetJWTSecretKey()))
}
