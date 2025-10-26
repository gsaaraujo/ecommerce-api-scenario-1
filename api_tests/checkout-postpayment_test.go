package apitests_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	testhelpers "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/test_helpers"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/stretchr/testify/suite"
)

type CheckoutPostpaymentSuite struct {
	suite.Suite
	customerDAO     daos.CustomerDAO
	cartDAO         daos.CartDAO
	cartItemDAO     daos.CartItemDAO
	inventoryDAO    daos.InventoryDAO
	addressDAO      daos.AddressDAO
	orderDAO        daos.OrderDAO
	orderItemDAO    daos.OrderItemDAO
	paymentDAO      daos.PaymentDAO
	productDAO      daos.ProductDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (c *CheckoutPostpaymentSuite) SetupSuite() {
	c.testEnvironment = testhelpers.NewTestEnvironment()
	err := c.testEnvironment.Start()
	c.Require().NoError(err)

	c.customerDAO = daos.NewCustomerDAO(c.testEnvironment.PgxPool())
	c.cartDAO = daos.NewCartDAO(c.testEnvironment.PgxPool())
	c.cartItemDAO = daos.NewCartItemDAO(c.testEnvironment.PgxPool())
	c.inventoryDAO = daos.NewInventoryDAO(c.testEnvironment.PgxPool())
	c.addressDAO = daos.NewAddressDAO(c.testEnvironment.PgxPool())
	c.orderDAO = daos.NewOrderDAO(c.testEnvironment.PgxPool())
	c.orderItemDAO = daos.NewOrderItemDAO(c.testEnvironment.PgxPool())
	c.paymentDAO = daos.NewPaymentDAO(c.testEnvironment.PgxPool())
	c.productDAO = daos.NewProductDAO(c.testEnvironment.PgxPool())
}

func (c *CheckoutPostpaymentSuite) SetupTest() {
	err := c.customerDAO.DeletAll()
	c.Require().NoError(err)

	err = c.addressDAO.DeletAll()
	c.Require().NoError(err)

	err = c.productDAO.DeletAll()
	c.Require().NoError(err)

	err = c.inventoryDAO.DeletAll()
	c.Require().NoError(err)

	err = c.cartDAO.DeletAll()
	c.Require().NoError(err)

	err = c.cartItemDAO.DeletAll()
	c.Require().NoError(err)

	err = c.orderDAO.DeletAll()
	c.Require().NoError(err)

	err = c.orderItemDAO.DeletAll()
	c.Require().NoError(err)

	err = c.paymentDAO.DeletAll()
	c.Require().NoError(err)
}

func (c *CheckoutPostpaymentSuite) Test1() {
	c.Run("when checking out, then it returns 200 and updates inventory, create a new order, order item and payment, and clears cart", func() {
		err := c.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		c.Require().NoError(err)
		err = c.addressDAO.Create(daos.AddressSchema{
			Id:          uuid.MustParse("9a6a0e64-4790-4ad2-99af-182f85bbac5b"),
			CustomerId:  uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			IsDefault:   true,
			Street:      "Maple Grove Lane",
			Number:      "4767",
			City:        "Austin",
			State:       "TX",
			ZipCode:     "78739",
			AddressLine: "4767 Maple Grove Lane, Austin, TX 78739",
			CreatedAt:   time.Now().UTC(),
		})
		c.Require().NoError(err)
		err = c.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		c.Require().NoError(err)
		err = c.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("7ab00199-6f9c-4af7-ad54-a02503226282"),
			Name:        "Kinesis Freestyle2 Wireless Ergonomic Keyboard",
			Description: utils.NewPointer("A split-design wireless ergonomic keyboard ..."),
			Price:       99286,
		})
		c.Require().NoError(err)
		err = c.inventoryDAO.Create(daos.InventorySchema{
			Id:            uuid.MustParse("cf23ee55-88c0-4898-ada4-15645c75645d"),
			ProductId:     uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			StockQuantity: 50,
			CreatedAt:     time.Now().UTC(),
		})
		c.Require().NoError(err)
		err = c.inventoryDAO.Create(daos.InventorySchema{
			Id:            uuid.MustParse("3fede283-d7f3-4423-bfe1-63163978c03f"),
			ProductId:     uuid.MustParse("7ab00199-6f9c-4af7-ad54-a02503226282"),
			StockQuantity: 50,
			CreatedAt:     time.Now().UTC(),
		})
		c.Require().NoError(err)
		err = c.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		c.Require().NoError(err)
		err = c.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("b999870f-f969-4d24-8955-499dbf3c689e"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Quantity:  8,
			CreatedAt: time.Now().UTC(),
		})
		c.Require().NoError(err)
		err = c.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("9052e9d7-84b0-4d6e-81aa-c59befb79088"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("7ab00199-6f9c-4af7-ad54-a02503226282"),
			Quantity:  4,
			CreatedAt: time.Now().UTC(),
		})
		c.Require().NoError(err)

		request, err := http.NewRequest("POST", c.testEnvironment.BaseUrl()+
			"/v1/checkout-postpayment?data_id=123456&customer_id=f59207c8-e837-4159-b67d-78c716510747&address_id=9a6a0e64-4790-4ad2-99af-182f85bbac5b", nil)
		c.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
		c.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := c.testEnvironment.Client().Do(request)
		c.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		c.Require().NoError(err)
		c.Equal(200, response.StatusCode)
		c.Equal("", string(body))

		orderSchema, err := c.orderDAO.FindOneByCustomerId(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		c.Require().NoError(err)
		c.Require().NotNil(orderSchema)
		c.Require().True(utils.IsValidUUID(orderSchema.CustomerId.String()))
		c.Require().Equal(int32(12), orderSchema.TotalQuantity)
		c.Require().Equal(int64(421136), orderSchema.TotalPrice)
		c.Require().WithinDuration(time.Now(), orderSchema.CreatedAt, 5*time.Second)

		orderItemSchema, err := c.orderItemDAO.FindAllByOrderId(orderSchema.Id)
		c.Require().NoError(err)
		c.Require().NotEmpty(orderItemSchema)
		c.Require().True(utils.IsValidUUID(orderItemSchema[0].Id.String()))
		c.Require().Equal(orderSchema.Id, orderItemSchema[0].OrderId)
		c.Require().Equal("c0981e5b-9cb7-4623-9713-55db0317dc1a", orderItemSchema[0].ProductId.String())
		c.Require().Equal(int32(8), orderItemSchema[0].Quantity)
		c.Require().Equal(int64(2999), orderItemSchema[0].Price)
		c.Require().WithinDuration(time.Now(), orderItemSchema[0].CreatedAt, 5*time.Second)

		c.Require().True(utils.IsValidUUID(orderItemSchema[1].Id.String()))
		c.Require().Equal(orderSchema.Id, orderItemSchema[1].OrderId)
		c.Require().Equal("7ab00199-6f9c-4af7-ad54-a02503226282", orderItemSchema[1].ProductId.String())
		c.Require().Equal(int32(4), orderItemSchema[1].Quantity)
		c.Require().Equal(int64(99286), orderItemSchema[1].Price)
		c.Require().WithinDuration(time.Now(), orderItemSchema[1].CreatedAt, 5*time.Second)

		paymentSchema, err := c.paymentDAO.FindOneByCustomerId(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		c.Require().NoError(err)
		c.Require().NotNil(paymentSchema)
		c.Require().True(utils.IsValidUUID(paymentSchema.OrderId.String()))
		c.Require().Equal("123456", paymentSchema.PaymentGatewayTransactionId)
		c.Require().Equal("mercado_pago", paymentSchema.PaymentGatewayName)
		c.Require().WithinDuration(time.Now(), paymentSchema.CreatedAt, 5*time.Second)

		cartItemSchema, err := c.cartItemDAO.FindAllByCartId(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"))
		c.Require().NoError(err)
		c.Require().Empty(cartItemSchema)
	})
}

func TestCheckoutPostpayment(t *testing.T) {
	suite.Run(t, new(CheckoutPostpaymentSuite))
}
