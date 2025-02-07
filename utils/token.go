package utils

import (
	"errors"
	
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(Email string, Role string) (string, error) {
	claims := JwtClaims{
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

func ValidateToken(SingnedToken string, c *gin.Context) (Claims *JwtClaims, err error) {
	token, err := jwt.ParseWithClaims(SingnedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRETE_KEY")), nil
		})
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		err = errors.New("token exprired")
		return
	}
	if claims.Email == " " {
		err=errors.New("invalid claims")
	}
	if claims.Issuer== " "{
		 err=errors.New("invalid claims")
	}
	c.Set("Issuer", claims.Issuer)
	
	return
}
