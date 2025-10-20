package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type AddAddressHandlerInput struct {
	City         any `validate:"required,string,notEmpty"`
	State        any `validate:"required,string,notEmpty"`
	ZipCode      any `validate:"required,string,notEmpty"`
	StreetName   any `validate:"required,string,notEmpty"`
	StreetNumber any `validate:"required,string,notEmpty"`
}

type AddAddressHandler struct {
	jsonBodyValidator webhttp.JSONBodyValidator
	addProductUsecase usecases.AddAddressUsecase
}

func NewAddAddressHandler(jsonBodyValidator webhttp.JSONBodyValidator, addProductUsecase usecases.AddAddressUsecase) AddAddressHandler {
	return AddAddressHandler{jsonBodyValidator, addProductUsecase}
}

func (a *AddAddressHandler) Handle(c echo.Context) error {
	var input AddAddressHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := a.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(*usecases.JwtAccessTokenClaims)

	err := a.addProductUsecase.Execute(usecases.AddAddressUsecaseInput{
		CustomerId:   uuid.MustParse(claims.Subject),
		City:         input.City.(string),
		State:        input.State.(string),
		ZipCode:      input.ZipCode.(string),
		StreetName:   input.StreetName.(string),
		StreetNumber: input.StreetNumber.(string),
	})
	if err == nil {
		return c.NoContent(204)
	}

	if err.Error() == "ZIP code does not match any location" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "ZIP code location does not match with provided city and state" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "state must be a valid 2-letter U.S. abbreviation (e.g. NY, CA)" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "ZIP code is invalid. It must be 5 digits (e.g. 12345)" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "street number must contain only digits (0-9)" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
