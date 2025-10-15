package apitests

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

type GetCartSuite struct {
	suite.Suite
	productDAO      daos.ProductDAO
	cartDAO         daos.CartDAO
	cartItemDAO     daos.CartItemDAO
	customerDAO     daos.CustomerDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (i *GetCartSuite) SetupSuite() {
	i.testEnvironment = testhelpers.NewTestEnvironment()
	err := i.testEnvironment.Start()
	i.Require().NoError(err)

	i.productDAO = daos.NewProductDAO(i.testEnvironment.PgxPool())
	i.cartDAO = daos.NewCartDAO(i.testEnvironment.PgxPool())
	i.cartItemDAO = daos.NewCartItemDAO(i.testEnvironment.PgxPool())
	i.customerDAO = daos.NewCustomerDAO(i.testEnvironment.PgxPool())
}

func (i *GetCartSuite) SetupTest() {
	err := i.customerDAO.DeletAll()
	i.Require().NoError(err)

	err = i.productDAO.DeletAll()
	i.Require().NoError(err)

	err = i.cartDAO.DeletAll()
	i.Require().NoError(err)

	err = i.cartItemDAO.DeletAll()
	i.Require().NoError(err)
}

func (i *GetCartSuite) Test1() {
	i.Run("given that there's products in cart, when getting cart, then returns 200", func() {
		err := i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		i.Require().NoError(err)
		err = i.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		i.Require().NoError(err)
		err = i.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("7ab00199-6f9c-4af7-ad54-a02503226282"),
			Name:        "Kinesis Freestyle2 Wireless Ergonomic Keyboard",
			Description: utils.NewPointer("A split-design wireless ergonomic keyboard ..."),
			Price:       99286,
		})
		i.Require().NoError(err)
		err = i.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("b0d11d8f-d8a7-4f2a-81ad-5df8c0db75d2"),
			Name:        "JBL Tune 520BT Wireless Headphones",
			Description: utils.NewPointer("Lightweight Bluetooth on-ear headphones ..."),
			Price:       22167,
		})
		i.Require().NoError(err)
		err = i.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		i.Require().NoError(err)
		err = i.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("b999870f-f969-4d24-8955-499dbf3c689e"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Quantity:  8,
			CreatedAt: time.Now().UTC(),
		})
		i.Require().NoError(err)
		err = i.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("9052e9d7-84b0-4d6e-81aa-c59befb79088"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("7ab00199-6f9c-4af7-ad54-a02503226282"),
			Quantity:  4,
			CreatedAt: time.Now().UTC(),
		})
		i.Require().NoError(err)
		err = i.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("1ff33790-7353-40c8-96bf-e7ab0bcacaa8"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("b0d11d8f-d8a7-4f2a-81ad-5df8c0db75d2"),
			Quantity:  6,
			CreatedAt: time.Now().UTC(),
		})
		i.Require().NoError(err)

		request, err := http.NewRequest("GET", i.testEnvironment.BaseUrl()+"/v1/cart", nil)
		i.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		i.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := i.testEnvironment.Client().Do(request)
		i.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		i.Require().NoError(err)
		i.Equal(200, response.StatusCode)
		i.JSONEq(`
			{
				"data": {
					"cartId": "bb8357b2-b978-4675-9521-ef2da0bd1747",
					"totalItems": 3,
					"totalQuantity": 18,
					"totalPrice": 554138,
					"items": [
						{
							"id": "b999870f-f969-4d24-8955-499dbf3c689e",
							"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
							"name": "ErgoClick Pro Wireless Mouse",
							"description": "Ergonomically designed wireless optical mouse ...",
							"quantity": 8,
							"price": 2999
						},					
						{
							"id": "9052e9d7-84b0-4d6e-81aa-c59befb79088",
							"productId": "7ab00199-6f9c-4af7-ad54-a02503226282",
							"name": "Kinesis Freestyle2 Wireless Ergonomic Keyboard",
							"description": "A split-design wireless ergonomic keyboard ...",
							"quantity": 4,
							"price": 99286
						},
						{
							"id": "1ff33790-7353-40c8-96bf-e7ab0bcacaa8",
							"productId": "b0d11d8f-d8a7-4f2a-81ad-5df8c0db75d2",
							"name": "JBL Tune 520BT Wireless Headphones",
							"description": "Lightweight Bluetooth on-ear headphones ...",
							"quantity": 6,
							"price": 22167
						}
					]
				}
			}
		`, string(body))
	})
}

func (i *GetCartSuite) Test2() {
	i.Run("given that the cart is empty, when getting cart, then returns 200", func() {
		err := i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		i.Require().NoError(err)
		err = i.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		i.Require().NoError(err)

		request, err := http.NewRequest("GET", i.testEnvironment.BaseUrl()+"/v1/cart", nil)
		i.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		i.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := i.testEnvironment.Client().Do(request)
		i.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		i.Require().NoError(err)
		i.Equal(200, response.StatusCode)
		i.JSONEq(`
			{
				"data": {
					"cartId": "bb8357b2-b978-4675-9521-ef2da0bd1747",
					"totalItems": 0,
					"totalQuantity": 0,
					"totalPrice": 0,
					"items": []
				}
			}
		`, string(body))
	})
}

func (i *GetCartSuite) Test3() {
	i.Run("given that there's no cart, when getting cart, then returns 409", func() {
		err := i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		i.Require().NoError(err)

		request, err := http.NewRequest("GET", i.testEnvironment.BaseUrl()+"/v1/cart", nil)
		i.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		i.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := i.testEnvironment.Client().Do(request)
		i.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		i.Require().NoError(err)
		i.Equal(409, response.StatusCode)
		i.JSONEq(`
			{
				"message": "cart not found"
			}
		`, string(body))
	})
}

func TestGetCartSuite(t *testing.T) {
	suite.Run(t, new(GetCartSuite))
}
