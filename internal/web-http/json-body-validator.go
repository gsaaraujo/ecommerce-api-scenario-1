package webhttp

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type JSONBodyValidator struct {
	validate *validator.Validate
}

func NewJSONBodyValidator() (JSONBodyValidator, error) {
	newValidator := validator.New(validator.WithRequiredStructEnabled())

	err := newValidator.RegisterValidation("string", isString)

	if err != nil {
		return JSONBodyValidator{}, err
	}

	err = newValidator.RegisterValidation("integer", isInteger)

	if err != nil {
		return JSONBodyValidator{}, err
	}

	err = newValidator.RegisterValidation("notEmpty", isNotEmpty)

	if err != nil {
		return JSONBodyValidator{}, err
	}

	err = newValidator.RegisterValidation("positive", isPositive)

	if err != nil {
		return JSONBodyValidator{}, err
	}

	err = newValidator.RegisterValidation("timeRFC3339", isTimeRFC3339)

	if err != nil {
		return JSONBodyValidator{}, err
	}

	err = newValidator.RegisterValidation("date", isDate)

	if err != nil {
		return JSONBodyValidator{}, err
	}

	jsonBodyValidator := JSONBodyValidator{
		validate: newValidator,
	}

	return jsonBodyValidator, nil
}

func isString(fieldLevel validator.FieldLevel) bool {
	return fieldLevel.Field().Kind() == reflect.String
}

func isInteger(fieldLevel validator.FieldLevel) bool {
	if fieldLevel.Field().Kind() != reflect.Float64 {
		return false
	}

	value := fieldLevel.Field().Float()
	return value == float64(int(value))
}

func isNotEmpty(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()

	if field.Kind() == reflect.String {
		return strings.TrimSpace(field.String()) != ""
	}

	return false
}

func isPositive(fieldLevel validator.FieldLevel) bool {
	if fieldLevel.Field().Kind() != reflect.Float64 {
		return false
	}

	value := fieldLevel.Field().Float()
	return value >= 0
}

func isTimeRFC3339(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()

	if field.Kind() != reflect.String {
		return false
	}

	_, err := time.Parse(time.RFC3339, field.String())
	return err == nil
}

func isDate(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()

	if field.Kind() != reflect.String {
		return false
	}

	_, err := time.Parse(time.DateOnly, field.String())
	return err == nil
}

func (j *JSONBodyValidator) Validate(body any) []string {
	err := j.validate.Struct(body)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := []string{}

		for _, validationError := range validationErrors {
			tag := validationError.Tag()
			field := strings.ToLower(validationError.Field()[:1]) + validationError.Field()[1:]

			switch tag {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", field))
			case "uuidv4":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be uuidv4", field))
			case "string":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be string", field))
			case "integer":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be integer", field))
			case "notEmpty":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must not be empty", field))
			case "positive":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be positive", field))
			case "timeRFC3339":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must follow format yyyy-mm-ddThh:mm:ssZ", field))
			case "date":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must follow format yyyy-mm-dd", field))
			}
		}

		return errorMessages
	}

	return []string{}
}
