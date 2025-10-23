package handlers

import (
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/labstack/echo/v4"
)

type CheckoutPostpaymentHandler struct {
	jsonBodyValidator          webhttp.JSONBodyValidator
	checkoutPostpaymentUsecase usecases.CheckoutPostpaymentUsecase
}

func NewCheckoutPostpaymentHandler(jsonBodyValidator webhttp.JSONBodyValidator, checkoutPostpaymentUsecase usecases.CheckoutPostpaymentUsecase) CheckoutPostpaymentHandler {
	return CheckoutPostpaymentHandler{jsonBodyValidator, checkoutPostpaymentUsecase}
}

func (a *CheckoutPostpaymentHandler) Handle(c echo.Context) error {
	dataId := c.QueryParam("data_id")
	customerId := c.QueryParam("customer_id")
	addressId := c.QueryParam("address_id")

	err := a.checkoutPostpaymentUsecase.Execute(usecases.CheckoutPostpaymentUsecaseInput{
		CustomerId:                  uuid.MustParse(customerId),
		AddressId:                   uuid.MustParse(addressId),
		PaymentGatewayTransactionId: dataId,
	})
	if err == nil {
		return c.NoContent(200)
	}

	return err
}
