package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = strings.ToLower(err.Field())
			element.Message = msgForTag(err.Tag(), err.Param())
			errors = append(errors, &element)
		}
	}
	return errors
}

func msgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value must be at least " + param + " characters"
	case "max":
		return "Value must be at most " + param + " characters"
	case "alphanum":
		return "Must be alphanumeric"
	case "oneof":
		return "Must be one of: " + strings.Join(strings.Split(param, " "), ", ")
	default:
		return "Invalid value"
	}
}
