package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type RemoveProductFromCartHandlerInput struct {
	ProductId any `validate:"required,uuid4"`
}

type RemoveProductFromCartHandler struct {
	jsonBodyValidator            webhttp.JSONBodyValidator
	removeProductFromCartUsecase usecases.RemoveProductFromCartUsecase
}

func NewRemoveProductFromCartHandler(
	jsonBodyValidation webhttp.JSONBodyValidator,
	removeProductFromCartUsecase usecases.RemoveProductFromCartUsecase,
) RemoveProductFromCartHandler {
	return RemoveProductFromCartHandler{jsonBodyValidation, removeProductFromCartUsecase}
}

func (r *RemoveProductFromCartHandler) Handle(c echo.Context) error {
	var input RemoveProductFromCartHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := r.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(*usecases.JwtAccessTokenClaims)

	err := r.removeProductFromCartUsecase.Execute(usecases.RemoveProductFromCartUsecaseInput{
		CustomerId: uuid.MustParse(claims.Subject),
		ProductId:  uuid.MustParse(input.ProductId.(string)),
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

	return err
}
