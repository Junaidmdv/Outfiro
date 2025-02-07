package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/smtp"
	"os"
	"outfiro/models"
	"time"
)

type OTP struct {	
	Value      string `json:"otp"`
	Email      string `json:"email" validate:"email"`	
}

func SendEmail(to []string, otp string) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)
	fmt.Println(os.Getenv("FROM_EMAIL"))
	fmt.Println(os.Getenv("FROM_EMAIL_PASSWORD"))

	fmt.Println(to[0])
	subject := "Subject: Your OTP Verification Code\n\n"
	body := fmt.Sprintf("Your One-Time Password (OTP) for verification is: [%s] \nThis code is valid for the next 5 minutes.\n\n\n\nPlease do not share this code with anyone for security reasons.", otp)
	message := subject + "\n" + body
	err := smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		to,
		[]byte(message),
	)

	return err
}

func (O *OTP) GenerateOTP(maxDigits uint32) error {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		return errors.New("error in the otp generation")
	}
	O.Value = fmt.Sprintf("%0*d", maxDigits, bi)
	return nil
}

func (O *OTP) Verifyotp(otp *models.OtpRecord) error {
	if time.Now().After(otp.ExpiryTime){
		return errors.New("otp exiprired,Please try again")
	}
    if O.Value !=otp.Value{
		return errors.New("otp is not matching.Please enter proper otp")
	}
	return nil
}
