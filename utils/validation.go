package utils

import (
	"fmt"
	"outfiro/models"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var regexAlphaNumSpace = regexp.MustCompile(`^[a-zA-Z\s]+$`)
var regexPhoneNum = regexp.MustCompile(`^(\+91|0)?[6-9][0-9]{9}$`)
var regexPincode = regexp.MustCompile(`^[1-9][0-9]{5}$`)
var regexPassword = regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)

func UserFormateError(errs validator.ValidationErrors) []string {
	var ErrorMessage []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("%s field is  required", err.Field()))
		case "email":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("Invlid %s.Please enter proper email address", err.Field()))
		case "alpha":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf(" %s cannot contain numbers. Please enter a valid name with letters only", err.Field()))
		case "min":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("%s must be 8 character long", err.Field()))
		case "max":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("%s must be 100 character long", err.Tag()))
		case "numeric":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("Invalid formate of %s,Enter proper value", err.Tag()))
		case "gt":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("%s Invalide formate of %s,Enter proper value", err.Tag(), err.Field()))
		case "alpha_space":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf(" Invalid formate %s,Enter proper value", err.Field()))
		case "eqfield":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("%s should be match with password,Please try again", err.Field()))
		case "url":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("Invalid %s.Enter proper %s", err.Field(), err.Field()))
		case "phone_number":
			ErrorMessage = append(ErrorMessage, "Invalid phone number. Please enter a valid number")
		case "pincode":
			ErrorMessage = append(ErrorMessage, "Invalid pincode.Please enter a valid pincode")
		case "password":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf(" Invalid %s.Enter valid %s ", err.Field(), err.Field()))
		case "lte":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("Invalid %s.Enter valid value", err.Field()))
		case "gte":
			ErrorMessage = append(ErrorMessage, fmt.Sprintf("%s Invalide formate of %s,Enter proper value", err.Field(), err.Field()))
		}
	}
	return ErrorMessage

}

func ValidateAlphaNumSpace(fl validator.FieldLevel) bool {
	return regexAlphaNumSpace.MatchString(fl.Field().String())
}
func ValidPhoneNum(fl validator.FieldLevel) bool {
	return regexPhoneNum.MatchString(fl.Field().String())
}
func ValidPincode(fl validator.FieldLevel) bool {
	return regexPincode.MatchString(fl.Field().String())
}
func ValidPassword(fl validator.FieldLevel) bool {
	return regexPassword.MatchString(fl.Field().String())
}

func ValidationOrderStatus(current_status string, newstatus string) error {
	//pending return
	OrderStatus := []string{models.Cancelled, models.Delivered, models.Shipped, models.Pending, models.Processing}
	var exist bool
	for _, items := range OrderStatus {
		if items == newstatus {
			exist = true
			break
		}
	}
	if !exist {
		return fmt.Errorf("invalid status code")
	}

	IsValid := map[string][]string{
		models.Pending:    {models.Cancelled, models.Processing},
		models.Processing: {models.Cancelled, models.Shipped},
		models.Shipped:    {models.Delivered},
		models.Delivered:  {models.Return},
	}
	//validate the cancelled product
	status, exist := IsValid[current_status]
	if !exist {
		panic("Invalid current status")
	}
	for _, items := range status {
		if items == newstatus {
			return nil
		}
	}
	return fmt.Errorf("%s can't be changed to %s",  current_status,newstatus)
}
