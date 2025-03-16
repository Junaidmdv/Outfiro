package utils

import (
	"math"
	"math/rand"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func RoundFloat(n float64) float64 {
	ratio := math.Pow(10, 2)
	return math.Round(n*ratio) / ratio
}

func GenerateReferralCode() string {

	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	code := make([]rune, 8)
	for i := range code {
		code[i] = letters[rnd.Intn(len(letters))]
	}
	return string(code)
}
