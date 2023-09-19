package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"unicode"
)

type CustomValidator struct {
	v *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.RegisterValidation("startsWithUpperCase", cv.validateStartsWithUpperCase)
	if err != nil {
		panic(err)
	}

	return cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.v.Struct(i)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]

		return cv.newValidationError(fieldErr)
	}
	return nil
}

func (cv *CustomValidator) newValidationError(fe validator.FieldError) error {
	switch fe.Tag() {
	case "required":
		return fmt.Errorf("field %s is required", fe.Field())
	case "email":
		return fmt.Errorf("field %s must be a valid email address", fe.Field())
	case "startsWithUpperCase":
		return fmt.Errorf("field %s must start with an upper case letter", fe.Field())
	case "min":
		return fmt.Errorf("field %s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Errorf("field %s must be at most %s characters", fe.Field(), fe.Param())
	case "gte":
		return fmt.Errorf("field %s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Errorf("field %s must be less than or equal to %s", fe.Field(), fe.Param())
	case "alpha":
		return fmt.Errorf("field %s must contain only alpha characters", fe.Field())
	case "oneof":
		return fmt.Errorf("field %s must be one of (%s)", fe.Field(), fe.Param())
	default:
		return fmt.Errorf("field %s is invalid", fe.Field())
	}
}

func (cv *CustomValidator) validateStartsWithUpperCase(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) == 0 {
		return true
	}
	firstChar := rune(value[0])
	return unicode.IsUpper(firstChar)
}
