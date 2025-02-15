package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(id uint, Email string, Role string) (string, error) {
	claims := JwtClaims{
		UserId: int(id),
		Email: Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Role,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 10)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	SignedToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRETE_KEY")))
	if err != nil {
		return "", err
	}

	return SignedToken, nil
}
