package apitests_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	testhelpers "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/test_helpers"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/stretchr/testify/suite"
)

type AddProductToCartSuite struct {
	suite.Suite
	productDAO      daos.ProductDAO
	cartDAO         daos.CartDAO
	cartItemDAO     daos.CartItemDAO
	customerDAO     daos.CustomerDAO
	inventoryDAO    daos.InventoryDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (a *AddProductToCartSuite) SetupSuite() {
	a.testEnvironment = testhelpers.NewTestEnvironment()
	err := a.testEnvironment.Start()
	a.Require().NoError(err)

	a.productDAO = daos.NewProductDAO(a.testEnvironment.PgxPool())
	a.cartDAO = daos.NewCartDAO(a.testEnvironment.PgxPool())
	a.cartItemDAO = daos.NewCartItemDAO(a.testEnvironment.PgxPool())
	a.customerDAO = daos.NewCustomerDAO(a.testEnvironment.PgxPool())
	a.inventoryDAO = daos.NewInventoryDAO(a.testEnvironment.PgxPool())
}

func (a *AddProductToCartSuite) SetupTest() {
	err := a.customerDAO.DeletAll()
	a.Require().NoError(err)

	err = a.productDAO.DeletAll()
	a.Require().NoError(err)

	err = a.cartDAO.DeletAll()
	a.Require().NoError(err)

	err = a.cartItemDAO.DeletAll()
	a.Require().NoError(err)

	err = a.inventoryDAO.DeletAll()
	a.Require().NoError(err)
}

func (a *AddProductToCartSuite) Test1() {
	a.Run("given that the product exists, when adding product to cart, then returns 204 and updates cart", func() {
		err := a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.inventoryDAO.Create(daos.InventorySchema{
			Id:            uuid.MustParse("cf23ee55-88c0-4898-ada4-15645c75645d"),
			ProductId:     uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			StockQuantity: 10,
			CreatedAt:     time.Now().UTC(),
		})
		a.Require().NoError(err)

		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 8
			}
		`))
		a.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := a.testEnvironment.Client().Do(request)
		a.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		a.Require().NoError(err)
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		cartItemSchema, err := a.cartItemDAO.FindAllByCartId(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"))
		a.Require().NoError(err)
		a.Require().Equal(1, len(cartItemSchema))
		a.Require().True(utils.IsValidUUID(cartItemSchema[0].Id.String()))
		a.Require().Equal(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"), cartItemSchema[0].CartId)
		a.Require().Equal("c0981e5b-9cb7-4623-9713-55db0317dc1a", cartItemSchema[0].ProductId.String())
		a.Require().Equal(int32(8), cartItemSchema[0].Quantity)
		a.Require().WithinDuration(time.Now(), cartItemSchema[0].CreatedAt, 5*time.Second)
	})
}

func (a *AddProductToCartSuite) Test2() {
	a.Run("given that the products exists, when adding product to cart multiple times, then returns 204 and updates cart", func() {
		err := a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("7ab00199-6f9c-4af7-ad54-a02503226282"),
			Name:        "Kinesis Freestyle2 Wireless Ergonomic Keyboard",
			Description: utils.NewPointer("A split-design wireless ergonomic keyboard ..."),
			Price:       99286,
		})
		a.Require().NoError(err)
		err = a.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.inventoryDAO.Create(daos.InventorySchema{
			Id:            uuid.MustParse("cf23ee55-88c0-4898-ada4-15645c75645d"),
			ProductId:     uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			StockQuantity: 10,
			CreatedAt:     time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.inventoryDAO.Create(daos.InventorySchema{
			Id:            uuid.MustParse("9c838f9b-2cdc-43f2-ad74-6e38a3a1bb89"),
			ProductId:     uuid.MustParse("7ab00199-6f9c-4af7-ad54-a02503226282"),
			StockQuantity: 10,
			CreatedAt:     time.Now().UTC(),
		})
		a.Require().NoError(err)

		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NoError(err)

		request1, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 8
			}
		`))
		a.Require().NoError(err)
		request1.Header.Add("Content-Type", "application/json")
		request1.Header.Add("Authorization", "Bearer "+accessToken)

		request2, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 4
			}
		`))
		a.Require().NoError(err)
		request2.Header.Add("Content-Type", "application/json")
		request2.Header.Add("Authorization", "Bearer "+accessToken)

		request3, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(`
			{
				"productId": "7ab00199-6f9c-4af7-ad54-a02503226282",
				"quantity": 2
			}
		`))
		a.Require().NoError(err)
		request3.Header.Add("Content-Type", "application/json")
		request3.Header.Add("Authorization", "Bearer "+accessToken)

		_, err = a.testEnvironment.Client().Do(request1)
		a.Require().NoError(err)
		_, err = a.testEnvironment.Client().Do(request2)
		a.Require().NoError(err)
		response, err := a.testEnvironment.Client().Do(request3)
		a.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		a.Require().NoError(err)
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		cartItemSchema, err := a.cartItemDAO.FindAllByCartId(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"))
		a.Require().NoError(err)
		a.Require().Equal(2, len(cartItemSchema))

		a.Require().True(utils.IsValidUUID(cartItemSchema[0].Id.String()))
		a.Require().Equal(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"), cartItemSchema[0].CartId)
		a.Require().Equal("c0981e5b-9cb7-4623-9713-55db0317dc1a", cartItemSchema[0].ProductId.String())
		a.Require().Equal(int32(12), cartItemSchema[0].Quantity)
		a.Require().WithinDuration(time.Now(), cartItemSchema[0].CreatedAt, 5*time.Second)

		a.Require().True(utils.IsValidUUID(cartItemSchema[1].Id.String()))
		a.Require().Equal(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"), cartItemSchema[1].CartId)
		a.Require().Equal("7ab00199-6f9c-4af7-ad54-a02503226282", cartItemSchema[1].ProductId.String())
		a.Require().Equal(int32(2), cartItemSchema[1].Quantity)
		a.Require().WithinDuration(time.Now(), cartItemSchema[1].CreatedAt, 5*time.Second)
	})
}

func (a *AddProductToCartSuite) Test3() {
	a.Run("given that the product exists, when adding product to cart and quantity is equals zero, then returns 409", func() {
		err := a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		a.Require().NoError(err)

		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 0
			}
		`))
		a.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := a.testEnvironment.Client().Do(request)
		a.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		a.Require().NoError(err)
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "product quantity cannot be zero"
			}
		`, string(body))
	})
}

func (a *AddProductToCartSuite) Test4() {
	a.Run("given that the product exists, when adding product to cart and quantity is higher than available in inventory, then returns 409", func() {
		err := a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		a.Require().NoError(err)
		err = a.inventoryDAO.Create(daos.InventorySchema{
			Id:            uuid.MustParse("cf23ee55-88c0-4898-ada4-15645c75645d"),
			ProductId:     uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			StockQuantity: 4,
			CreatedAt:     time.Now().UTC(),
		})
		a.Require().NoError(err)

		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 8
			}
		`))
		a.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := a.testEnvironment.Client().Do(request)
		a.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		a.Require().NoError(err)
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "product quantity exceeds the stock available"
			}
		`, string(body))
	})
}

func (a *AddProductToCartSuite) Test5() {
	a.Run("given that the product does not exist, when adding product to cart, then returns 409", func() {
		err := a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		a.Require().NoError(err)

		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 8
			}
		`))
		a.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := a.testEnvironment.Client().Do(request)
		a.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		a.Require().NoError(err)
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "product not found"
			}
		`, string(body))
	})
}

func (a *AddProductToCartSuite) Test6() {
	a.Run("when adding product to cart and body is invalid, then returns 400", func() {
		templates := []map[string]string{
			{
				"body": `{}`,
				"error": `[
					"productId is required",
					"quantity is required"
				]`,
			},
			{
				"body": `{
					"productId": null,
					"quantity": null
				}`,
				"error": `[
					"productId is required",
					"quantity is required"
				]`,
			},
			{
				"body": `{
					"productId": "",
					"quantity": ""
				}`,
				"error": `[
					"productId must be uuidv4",
					"quantity must be integer"
				]`,
			},
			{
				"body": `{
					"productId": " ",
					"quantity": " "
				}`,
				"error": `[
					"productId must be uuidv4",
					"quantity must be integer"
				]`,
			},
			{
				"body": `{
					"productId": 1,
					"quantity": 1
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": 1.5,
					"quantity": 1.5
				}`,
				"error": `[
					"productId must be uuidv4",
					"quantity must be integer"
				]`,
			},
			{
				"body": `{
					"productId": -1,
					"quantity": -1
				}`,
				"error": `[
					"productId must be uuidv4",
					"quantity must be positive"
				]`,
			},
			{
				"body": `{
					"productId": true,
					"quantity": false
				}`,
				"error": `[
					"productId must be uuidv4",
					"quantity must be integer"
				]`,
			},
			{
				"body": `{
					"productId": {},
					"quantity": {}
				}`,
				"error": `[
					"productId must be uuidv4",
					"quantity must be integer"
				]`,
			},
			{
				"body": `{
					"productId": [],
					"quantity": []
				}`,
				"error": `[
					"productId must be uuidv4",
					"quantity must be integer"
				]`,
			},
		}

		for _, template := range templates {
			request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-product-to-cart", strings.NewReader(template["body"]))
			a.Require().NoError(err)
			accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
			a.Require().NoError(err)
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", "Bearer "+accessToken)

			response, err := a.testEnvironment.Client().Do(request)
			a.Require().NoError(err)

			body, err := io.ReadAll(response.Body)
			a.Require().NoError(err)

			a.Equal(400, response.StatusCode)
			a.JSONEq(fmt.Sprintf(`
				{
					"message": %s
				}
			`, template["error"]), string(body))
		}
	})
}

func TestAddProductToCart(t *testing.T) {
	suite.Run(t, new(AddProductToCartSuite))
}
