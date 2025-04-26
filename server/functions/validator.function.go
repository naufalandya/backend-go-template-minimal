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

var validationMessages = map[string]string{
	"required":     "%s is required~!",
	"min":          "%s must be at least %s character(s)~",
	"max":          "%s must be at most %s character(s)~",
	"len":          "%s must be exactly %s character(s)~",
	"email":        "%s must be a valid email~",
	"gt":           "%s must be greater than %s~",
	"gte":          "%s must be greater than or equal to %s~",
	"lt":           "%s must be less than %s~",
	"lte":          "%s must be less than or equal to %s~",
	"eq":           "%s must be equal to %s~",
	"ne":           "%s must not be equal to %s~",
	"oneof":        "%s must be one of [%s]~",
	"url":          "%s must be a valid URL~",
	"uuid":         "%s must be a valid UUID~",
	"alphanum":     "%s must contain only letters and numbers~",
	"numeric":      "%s must be a number~",
	"contains":     "%s must contain '%s'~",
	"startswith":   "%s must start with '%s'~",
	"endswith":     "%s must end with '%s'~",
	"excludes":     "%s must not contain '%s'~",
	"ip":           "%s must be a valid IP address~",
	"ipv4":         "%s must be a valid IPv4 address~",
	"ipv6":         "%s must be a valid IPv6 address~",
	"datetime":     "%s must be a valid datetime (format: %s)~",
	"boolean":      "%s must be true or false~",
	"isdefault":    "%s must have the default value~",
	"base64":       "%s must be a valid base64 string~",
	"lowercase":    "%s must be lowercase~",
	"uppercase":    "%s must be uppercase~",
	"ascii":        "%s must contain only ascii characters~",
	"printascii":   "%s must contain only printable ascii characters~",
	"containsany":  "%s must contain at least one of these characters: %s~",
	"excludesall":  "%s must not contain any of these characters: %s~",
	"excludesrune": "%s must not contain the rune '%s'~",
	// ✨ Add more if needed nya~!
}

func ValidateStruct(data interface{}) []*ValidationError {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var errors []*ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		tag := e.Tag()
		messageTemplate, ok := validationMessages[tag]
		var message string

		if ok {
			if e.Param() != "" {
				message = fmt.Sprintf(messageTemplate, e.Field(), e.Param())
			} else {
				message = fmt.Sprintf(messageTemplate, e.Field())
			}
		} else {
			message = fmt.Sprintf("%s is not valid~", e.Field())
		}

		errors = append(errors, &ValidationError{
			Field:   e.Field(),
			Message: message,
		})
	}

	return errors
}
