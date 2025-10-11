package handlers

import (
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type PublishProductHandlerInput struct {
	ProductId any `validate:"required,uuid4"`
}

type PublishProductHandler struct {
	jsonBodyValidator     webhttp.JSONBodyValidator
	publishProductUsecase usecases.PublishProductUsecase
}

func NewPublishProductHandler(jsonBodyValidator webhttp.JSONBodyValidator, publishProductUsecase usecases.PublishProductUsecase) PublishProductHandler {
	return PublishProductHandler{jsonBodyValidator, publishProductUsecase}
}

func (p *PublishProductHandler) Handle(c echo.Context) error {
	var input PublishProductHandlerInput

	if err := c.Bind(&input); err != nil {
		return c.NoContent(415)
	}

	if messages := p.jsonBodyValidator.Validate(input); len(messages) > 0 {
		return c.JSON(400, map[string]any{"message": messages})
	}

	err := p.publishProductUsecase.Execute(usecases.PublishProductUsecaseInput{
		ProductId: uuid.MustParse(input.ProductId.(string)),
	})
	if err == nil {
		return c.NoContent(204)
	}

	if err.Error() == "product not found" {
		return c.JSON(409, map[string]any{"message": err.Error()})
	}

	return err
}
