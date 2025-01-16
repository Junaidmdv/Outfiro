package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/smtp"
	"os"
	"time"

	"gorm.io/gorm"
)
type OTP struct {
	gorm.Model
	Value       string    `json:"otp" validate:"required" `
	Email     string       `json:"email" validate:"required,email"`
}


func SendEmail(to []string,otp string)error{
    auth:= smtp.PlainAuth(
		"",
		os.Getenv("FROM EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)
	subject:="Subject: Your OTP Verification Code"
    body:=fmt.Sprintf("Your One-Time Password (OTP) for verification is: [%s] \n This code is valid for the next 5 minutes. Please do not share this code with anyone for security reasons.",otp)
	message:=subject+"\n"+body
	err:=smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		to,
		[]byte(message),
	)

	return err
}

func (O *OTP) GenerateOTP(maxDigits uint32)error{
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		return errors.New("error in the otp generation")
	}
	O.Value= fmt.Sprintf("%0*d", maxDigits, bi)
	 return nil
}

func (O *OTP)Verifyotp(otp string)error{
  otpExpirationtime:=5*time.Minute

   if time.Now().After(O.CreatedAt.Add(otpExpirationtime)){
     return errors.New("the otp expired")
   }
   if O.Value !=otp{
      return errors.New("invalid otp")
   }

   return nil
}