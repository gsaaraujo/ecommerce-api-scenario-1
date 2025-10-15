package handlers

import (
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type SignUpHandlerInput struct {
	Name     any `validate:"required,string,notEmpty"`
	Email    any `validate:"required,string,notEmpty"`
	Password any `validate:"required,string,notEmpty"`
}

type SignUpHandler struct {
	jsonBodyValidator webhttp.JSONBodyValidator
	signUpUsecase     usecases.SignUpUsecase
}

func NewSignUpHandler(jsonBodyValidator webhttp.JSONBodyValidator, signUpUsecase usecases.SignUpUsecase) SignUpHandler {
	return SignUpHandler{jsonBodyValidator, signUpUsecase}
}

func (r SignUpHandler) Handle(c echo.Context) error {
	var input SignUpHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := r.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	err := r.signUpUsecase.Execute(usecases.SignUpUsecaseInput{
		Name:     input.Name.(string),
		Email:    input.Email.(string),
		Password: input.Password.(string),
	})
	if err == nil {
		return c.NoContent(204)
	}

	if err.Error() == "name must be at least 2 characters" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "email address is invalid" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "password must be at least 6 characters" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "this email address has already been taken by someone" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
