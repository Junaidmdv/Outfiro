package models

import (
	"strings"
	"github.com/go-playground/validator/v10"
)


type OtpData struct{
	Email string `json:"email,omitempty" validate:"required"`
} 


type OtpVerify struct{
	User *OtpData `json:"user,omitempty" validate:"required"`
	OtpValue   string `json:"code,omitempty"  validate:"required"`
}

func ContainsSpecialChars(fl validator.FieldLevel)bool{
   specialcharactar:="!@#$%^&*"
   for _,val:=range specialcharactar {
       strings.Contains(fl.Field().String(),string(val))
	   return true
   }
   return false
 
}