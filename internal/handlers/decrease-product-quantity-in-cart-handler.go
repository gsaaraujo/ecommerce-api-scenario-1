package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type DecreaseProductQuantityInCartHandlerInput struct {
	ProductId any `validate:"required,uuid4"`
	Quantity  any `validate:"required,integer,positive"`
}

type DecreaseProductQuantityInCartHandler struct {
	jsonBodyValidator                    webhttp.JSONBodyValidator
	decreaseProductQuantityInCartUsecase usecases.DecreaseProductQuantityInCartUsecase
}

func NewDecreaseProductQuantityInCartHandler(
	jsonBodyValidation webhttp.JSONBodyValidator,
	decreaseProductQuantityInCartUsecase usecases.DecreaseProductQuantityInCartUsecase,
) DecreaseProductQuantityInCartHandler {
	return DecreaseProductQuantityInCartHandler{jsonBodyValidation, decreaseProductQuantityInCartUsecase}
}

func (r *DecreaseProductQuantityInCartHandler) Handle(c echo.Context) error {
	var input DecreaseProductQuantityInCartHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := r.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(*usecases.JwtAccessTokenClaims)

	err := r.decreaseProductQuantityInCartUsecase.Execute(usecases.DecreaseProductQuantityInCartUsecaseInput{
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

	if err.Error() == "you cannot decrease the quantity of product with a value equal to zero" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
