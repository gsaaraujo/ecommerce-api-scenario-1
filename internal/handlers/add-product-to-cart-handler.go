package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type AddProductToCartHandlerInput struct {
	ProductId any `validate:"required,uuid4"`
	Quantity  any `validate:"required,integer,positive"`
}

type AddProductToCartHandler struct {
	jsonBodyValidator       webhttp.JSONBodyValidator
	addProductToCartUsecase usecases.AddProductToCartUsecase
}

func NewAddProductToCartHandler(jsonBodyValidator webhttp.JSONBodyValidator, addProductToCartUsecase usecases.AddProductToCartUsecase) AddProductToCartHandler {
	return AddProductToCartHandler{jsonBodyValidator, addProductToCartUsecase}
}

func (a *AddProductToCartHandler) Handle(c echo.Context) error {
	var input AddProductToCartHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := a.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(*usecases.JwtAccessTokenClaims)

	err := a.addProductToCartUsecase.Execute(usecases.AddProductToCartUsecaseInput{
		CustomerId: uuid.MustParse(claims.Subject),
		ProductId:  uuid.MustParse(input.ProductId.(string)),
		Quantity:   int32(input.Quantity.(float64)),
	})
	if err == nil {
		return c.NoContent(204)
	}

	if err.Error() == "product not found" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "product quantity cannot be zero" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "product quantity exceeds the stock available" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
