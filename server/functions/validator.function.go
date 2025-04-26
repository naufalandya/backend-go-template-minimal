package functions

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(data interface{}) []*ValidationError {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var errors []*ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		var message string

		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required~!", e.Field())
		case "min":
			message = fmt.Sprintf("%s must be at least %s character(s)/value~", e.Field(), e.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s character(s)/value~", e.Field(), e.Param())
		case "email":
			message = fmt.Sprintf("%s must be a valid email~", e.Field())
		case "gt":
			message = fmt.Sprintf("%s must be greater than %s~", e.Field(), e.Param())
		case "lt":
			message = fmt.Sprintf("%s must be less than %s~", e.Field(), e.Param())
		default:
			message = fmt.Sprintf("%s is not valid~", e.Field())
		}

		errors = append(errors, &ValidationError{
			Field:   e.Field(),
			Message: message,
		})
	}

	return errors
}
