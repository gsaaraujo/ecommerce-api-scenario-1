package handlers

import (
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type AddProductHandlerInput struct {
	Name        any `validate:"required,string,notEmpty"`
	Description any `validate:"omitempty,string,notEmpty"`
	Price       any `validate:"required,integer,positive"`
}

type AddProductHandler struct {
	jsonBodyValidator webhttp.JSONBodyValidator
	addProductUsecase usecases.AddProductUsecase
}

func NewAddProductHandler(jsonBodyValidator webhttp.JSONBodyValidator, addProductUsecase usecases.AddProductUsecase) AddProductHandler {
	return AddProductHandler{jsonBodyValidator, addProductUsecase}
}

func (a *AddProductHandler) Handle(c echo.Context) error {
	var input AddProductHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := a.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	var description *string = nil

	if input.Description != nil {
		s := input.Description.(string)
		description = &s
	}

	err := a.addProductUsecase.Execute(usecases.AddProductUsecaseInput{
		Name:        input.Name.(string),
		Description: description,
		Price:       int64(input.Price.(float64)),
	})
	if err == nil {
		return c.NoContent(204)
	}

	if err.Error() == "the product price cannot be zero" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
