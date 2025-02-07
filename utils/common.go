package utils
import(
	"regexp"
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

func VerifyPassword(hash string,password string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err
}