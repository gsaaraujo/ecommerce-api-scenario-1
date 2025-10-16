package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type IncreaseProductQuantityInCartHandlerInput struct {
	ProductId any `validate:"required,uuid4"`
	Quantity  any `validate:"required,integer,positive"`
}

type IncreaseProductQuantityInCartHandler struct {
	jsonBodyValidator                    webhttp.JSONBodyValidator
	increaseProductQuantityInCartUsecase usecases.IncreaseProductQuantityInCartUsecase
}

func NewIncreaseProductQuantityInCartHandler(
	jsonBodyValidation webhttp.JSONBodyValidator,
	increaseProductQuantityInCartUsecase usecases.IncreaseProductQuantityInCartUsecase,
) IncreaseProductQuantityInCartHandler {
	return IncreaseProductQuantityInCartHandler{jsonBodyValidation, increaseProductQuantityInCartUsecase}
}

func (r *IncreaseProductQuantityInCartHandler) Handle(c echo.Context) error {
	var input IncreaseProductQuantityInCartHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := r.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(*usecases.JwtAccessTokenClaims)

	err := r.increaseProductQuantityInCartUsecase.Execute(usecases.IncreaseProductQuantityInCartUsecaseInput{
		CustomerId: uuid.MustParse(claims.Subject),
		ProductId:  uuid.MustParse(input.ProductId.(string)),
		Quantity:   int32(input.Quantity.(float64)),
	})
	if err == nil {
		return c.NoContent(204)
	}

	if err.Error() == "cart not found" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "product not found in cart" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "you cannot increase the quantity of product with a value equal to zero" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	if err.Error() == "product quantity exceeds the stock available" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
