package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var regexAlphaNumSpace = regexp.MustCompile(`^[a-zA-Z\s]+$`)

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
		}
	}
	return ErrorMessage

}

func ValidateAlphaNumSpace(fl validator.FieldLevel) bool {
	return regexAlphaNumSpace.MatchString(fl.Field().String())
}
