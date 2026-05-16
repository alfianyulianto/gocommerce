package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type errorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   any    `json:"value"`
}

var Validate *validator.Validate = validator.New()

func pascalCaseToReadable(s string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	result := re.ReplaceAllString(s, `$1 $2`)
	return result
}

func pascalToSnakeCase(s string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)
	return strings.ToLower(snake)
}

func ParseMessage(err error) map[string]errorMessage {
	errors := make(map[string]errorMessage)

	for _, e := range err.(validator.ValidationErrors) {
		errors[pascalToSnakeCase(e.Field())] = errorMessage{
			Field:   pascalToSnakeCase(e.Field()),
			Message: getCustomeMessage(e),
			Tag:     e.Tag(),
			Value:   e.Value(),
		}
	}

	return errors
}

func getCustomeMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "email":
		return fmt.Sprintf("%s must be a valid email address", pascalCaseToReadable(e.Field()))
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", pascalCaseToReadable(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s must not be at most %s characters", pascalCaseToReadable(e.Field()), e.Param())
	case "number":
		return fmt.Sprintf("%s must be valid number", pascalCaseToReadable(e.Field()))
	case "oneof":
		return fmt.Sprintf("%s must be one of '%s'.", pascalCaseToReadable(e.Field()), e.Param())
	case "required":
		return fmt.Sprintf("%s is required", pascalCaseToReadable(e.Field()))
	default:
		return e.Error()
	}
}
