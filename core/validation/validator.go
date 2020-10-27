// Core package for validating JSON input. Some useful links that was helpful:
// # Registering custom validator ->
// https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835
// # Custom error message ->
// https://learning-cloud-native-go.github.io/docs/a12.adding_validator.v9/

package validation


import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)


// A new `*validator.Validate` with the custom validators setup with it.
func NewValidator() *validator.Validate {
	validate := validator.New()
	RegisterTagName(validate)
	RegisterCustomValidators(validate)
	return validate
}


// This function helps to get the JSON name instead of the struct attribute name
func RegisterTagName(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
        name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })
}


// Registers custom validators
func RegisterCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation("uppercase", func(field validator.FieldLevel) bool {
		return HasUppercase(field.Field().String())
	})
	validate.RegisterValidation("lowercase", func(field validator.FieldLevel) bool {
		return HasLowercase(field.Field().String())
	})
	validate.RegisterValidation("digit", func(field validator.FieldLevel) bool {
		return HasDigit(field.Field().String())
	})
	validate.RegisterValidation("punctuation", func(field validator.FieldLevel) bool {
		return HasPunctuation(field.Field().String())
	})
}


// Custom error messages for different validators (both builtin and custom)
func ToCustomErrorMessage(err error) map[string]string {
	var errorMessage map[string]string = make(map[string]string)
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
        for _, err := range fieldErrors {
            switch err.Tag() {
			case "required":
				errorMessage[err.Field()] = fmt.Sprintf("%s is a required field.", err.Field())
			case "uppercase":
				errorMessage[err.Field()] = fmt.Sprintf("%s must contain at least one capital letter.", err.Field())
			case "lowercase":
				errorMessage[err.Field()] = fmt.Sprintf("%s must contain at least one small letter.", err.Field())
			case "digit":
				errorMessage[err.Field()] = fmt.Sprintf("%s must contain at least one digit.", err.Field())
			case "punctuation":
				errorMessage[err.Field()] = fmt.Sprintf("%s must contain at least one punctuation.", err.Field())
			case "email":
				errorMessage[err.Field()] = fmt.Sprintf("%s must be a valid email.", err.Field())
			case "min":
				errorMessage[err.Field()] = fmt.Sprintf("%s must be of length at least %s.", err.Field(), err.Param())
			case "max":
				errorMessage[err.Field()] = fmt.Sprintf("%s must be of length at best %d.", err.Field(), err.Param())
			default:
				errorMessage[err.Field()] = fmt.Sprintf("Missing custom error message for %s on tag %s.", err.Field(), err.Tag())
            }
        }
    }
	return errorMessage
}
