package utils

import (
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId   uint
	Email    string
	UserRole string
	jwt.RegisteredClaims
}

var Secretekey = []byte(os.Getenv("TOKEN_SECRETE_KEY"))

func CreateToken(userid uint, email string, userRole string) (string, error) {
	claims := Claims{
		UserId:   userid,
		Email:    email,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SignedToken, err := token.SignedString(Secretekey)
	if err != nil {
		return "", err
	}

	return SignedToken, nil

}

func ValidateToken(Claims *Claims) error {
	if !Claims.ExpiresAt.IsZero() && Claims.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("tokens are expired")
	}
	if Claims.UserId == 0 {
		return fmt.Errorf("user id is required")
	}
	if Claims.Email == " " {
		return fmt.Errorf("user email is required ")
	}
	if Claims.UserRole == " " {
		return fmt.Errorf("user role is required")
	}
	return nil
}
