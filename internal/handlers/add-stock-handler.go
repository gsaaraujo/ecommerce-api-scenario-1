package handlers

import (
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type AddStockHandlerInput struct {
	InventoryId any `validate:"required,uuid4"`
	Stock       any `validate:"required,integer,positive"`
}

type AddStockHandler struct {
	jsonBodyValidator webhttp.JSONBodyValidator
	addStockUsecase   usecases.AddStockUsecase
}

func NewAddStockHandler(jsonBodyValidator webhttp.JSONBodyValidator, addStockUsecase usecases.AddStockUsecase) AddStockHandler {
	return AddStockHandler{jsonBodyValidator, addStockUsecase}
}

func (a *AddStockHandler) Handle(c echo.Context) error {
	var input AddStockHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := a.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	err := a.addStockUsecase.Execute(usecases.AddStockUsecaseInput{
		InventoryId: uuid.MustParse(input.InventoryId.(string)),
		Stock:       int32(input.Stock.(float64)),
	})
	if err == nil {
		return c.NoContent(204)
	}

	if err.Error() == "inventory not found" {
		return c.JSON(404, map[string]any{"message": err.Error()})
	}

	if err.Error() == "stock quantity must be higher than zero" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
