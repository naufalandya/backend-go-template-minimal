package functions

import (
	"errors"
	"html"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func SuperSecureSanitize(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return sanitizeString(v)
	case int, int8, int16, int32, int64:
		return v, nil
	case float32, float64:
		return v, nil
	case bool:
		return v, nil
	case time.Time:
		if v.Year() < 1900 || v.Year() > 2100 {
			return nil, errors.New("date is suspicious~ (๑•﹏•)")
		}
		return v, nil
	case []string:
		var sanitized []string
		for _, str := range v {
			s, err := sanitizeString(str)
			if err != nil {
				return nil, err
			}
			sanitized = append(sanitized, s)
		}
		return sanitized, nil
	default:
		return nil, errors.New("unsupported input type detected~ (｡•́︿•̀｡)")
	}
}

func sanitizeString(value string) (string, error) {
	value = strings.TrimSpace(value)
	value = html.EscapeString(value)

	injectionPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\b(SELECT|INSERT|DELETE|DROP|UPDATE|EXEC|UNION|--|;|')\b`),
		regexp.MustCompile(`(?i)<script.*?>.*?</script>`),
		regexp.MustCompile(`[&|;$%@"<>(){}[\]]`),
		regexp.MustCompile(`\.\./|\.\.\\`),
		regexp.MustCompile(`(?i)(eval|alert|onerror|onload|prompt)\s*\(`),
	}

	for _, pattern := range injectionPatterns {
		if pattern.MatchString(value) {
			return "", errors.New("suspicious string detected~ (；⌣̀_⌣́)")
		}
	}

	value = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) && !unicode.IsSymbol(r) {
			return r
		}
		return -1
	}, value)

	return value, nil
}

func AutoSanitizeStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("need a non-nil pointer to a struct~! (`･ω･´)")
	}
	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return errors.New("input is not a struct pointer~ (；⌣̀_⌣́)")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Struct && fieldType.Type != reflect.TypeOf(time.Time{}) {
			err := AutoSanitizeStruct(field.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		switch field.Kind() {
		case reflect.String:
			sanitized, err := sanitizeString(field.String())
			if err != nil {
				return err
			}
			field.SetString(sanitized)

		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				sanitized := reflect.MakeSlice(field.Type(), field.Len(), field.Cap())
				for j := 0; j < field.Len(); j++ {
					cleaned, err := sanitizeString(field.Index(j).String())
					if err != nil {
						return err
					}
					sanitized.Index(j).SetString(cleaned)
				}
				field.Set(sanitized)
			}

		case reflect.Struct:
			if field.Type() == reflect.TypeOf(time.Time{}) {
				if field.Interface().(time.Time).Year() < 1900 || field.Interface().(time.Time).Year() > 2100 {
					return errors.New("suspicious date detected~ (｡•́︿•̀｡)")
				}
			}

		case reflect.Int, reflect.Int64, reflect.Float64, reflect.Bool:

		default:
			continue
		}
	}

	return nil
}

func AutoSuperSanitizeStruct(input interface{}) error {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("input must be a non-nil pointer to a struct~ (｡•́︿•̀｡)")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("input must point to a struct~")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			sanitized, err := SuperSecureSanitize(field.Interface())
			if err != nil {
				return err
			}
			field.SetString(sanitized.(string))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sanitized, err := SuperSecureSanitize(field.Interface())
			if err != nil {
				return err
			}
			field.SetInt(reflect.ValueOf(sanitized).Int())

		case reflect.Float32, reflect.Float64:
			sanitized, err := SuperSecureSanitize(field.Interface())
			if err != nil {
				return err
			}
			field.SetFloat(reflect.ValueOf(sanitized).Float())

		case reflect.Bool:
			sanitized, err := SuperSecureSanitize(field.Interface())
			if err != nil {
				return err
			}
			field.SetBool(reflect.ValueOf(sanitized).Bool())

		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				sanitized, err := SuperSecureSanitize(field.Interface())
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(sanitized))
			}

		case reflect.Struct:
			if field.Type().String() == "time.Time" {
				sanitized, err := SuperSecureSanitize(field.Interface())
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(sanitized))
			}
		}
	}

	return nil
}
