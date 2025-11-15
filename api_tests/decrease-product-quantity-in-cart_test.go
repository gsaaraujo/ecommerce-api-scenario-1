package apitests

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

type DecreaseProductQuantityInCartSuite struct {
	suite.Suite
	productDAO      daos.ProductDAO
	cartDAO         daos.CartDAO
	cartItemDAO     daos.CartItemDAO
	customerDAO     daos.CustomerDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (i *DecreaseProductQuantityInCartSuite) SetupSuite() {
	i.testEnvironment = testhelpers.NewTestEnvironment()
	i.testEnvironment.Start()

	i.productDAO = daos.NewProductDAO(i.testEnvironment.PgxPool())
	i.cartDAO = daos.NewCartDAO(i.testEnvironment.PgxPool())
	i.cartItemDAO = daos.NewCartItemDAO(i.testEnvironment.PgxPool())
	i.customerDAO = daos.NewCustomerDAO(i.testEnvironment.PgxPool())
}

func (i *DecreaseProductQuantityInCartSuite) SetupTest() {
	i.customerDAO.DeletAll()
	i.productDAO.DeletAll()
	i.cartDAO.DeletAll()
	i.cartItemDAO.DeletAll()
}

func (i *DecreaseProductQuantityInCartSuite) Test1() {
	i.Run("given that the product is in cart, when decreasing product quantity, then returns 204 and updates cart", func() {
		i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		i.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		i.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		i.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("b999870f-f969-4d24-8955-499dbf3c689e"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Quantity:  8,
			CreatedAt: time.Now().UTC(),
		})

		request := utils.GetOrThrow(http.NewRequest("POST", i.testEnvironment.BaseUrl()+"/v1/decrease-product-quantity-in-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 2
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(i.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		i.Equal(204, response.StatusCode)
		i.Equal("", string(body))

		cartItemSchema := i.cartItemDAO.FindAllByCartId(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"))
		i.Require().Equal(int32(6), cartItemSchema[0].Quantity)
	})
}

func (i *DecreaseProductQuantityInCartSuite) Test2() {
	i.Run("given that the product is in cart, when decreasing product quantity and it reaches zero, then returns 204 and deletes item from cart", func() {
		i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		i.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		i.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		i.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("b999870f-f969-4d24-8955-499dbf3c689e"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Quantity:  8,
			CreatedAt: time.Now().UTC(),
		})

		request := utils.GetOrThrow(http.NewRequest("POST", i.testEnvironment.BaseUrl()+"/v1/decrease-product-quantity-in-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 8
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(i.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		i.Equal(204, response.StatusCode)
		i.Equal("", string(body))

		cartItemSchema := i.cartItemDAO.FindAllByCartId(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"))
		i.Require().Empty(cartItemSchema)
	})
}

func (i *DecreaseProductQuantityInCartSuite) Test3() {
	i.Run(`given that the product is in cart, when decreasing product quantity and quantity is higher than current item quantity, 
	then returns 204 and deletes item from cart`, func() {
		i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		i.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		i.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		i.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("b999870f-f969-4d24-8955-499dbf3c689e"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Quantity:  2,
			CreatedAt: time.Now().UTC(),
		})

		request := utils.GetOrThrow(http.NewRequest("POST", i.testEnvironment.BaseUrl()+"/v1/decrease-product-quantity-in-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 8
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(i.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		i.Equal(204, response.StatusCode)
		i.Equal("", string(body))

		cartItemSchema := i.cartItemDAO.FindAllByCartId(uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"))
		i.Require().Empty(cartItemSchema)
	})
}

func (i *DecreaseProductQuantityInCartSuite) Test4() {
	i.Run("given that the product is not in cart, when decreasing product quantity, then returns 409", func() {
		i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		i.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		i.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})

		request := utils.GetOrThrow(http.NewRequest("POST", i.testEnvironment.BaseUrl()+"/v1/decrease-product-quantity-in-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 2
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(i.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		i.Equal(409, response.StatusCode)
		i.JSONEq(`
			{
				"message": "product not found in cart"
			}
		`, string(body))
	})
}

func (i *DecreaseProductQuantityInCartSuite) Test5() {
	i.Run("when decreasing product quantity and quantity is zero, then returns 409", func() {
		i.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})

		request := utils.GetOrThrow(http.NewRequest("POST", i.testEnvironment.BaseUrl()+"/v1/decrease-product-quantity-in-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a",
				"quantity": 0
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(i.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		i.Equal(409, response.StatusCode)
		i.JSONEq(`
			{
				"message": "you cannot decrease the quantity of product with a value equal to zero"
			}
		`, string(body))
	})
}

func (i *DecreaseProductQuantityInCartSuite) Test6() {
	i.Run("when decreasing product quantity and body is invalid, then returns 400", func() {
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
			request := utils.GetOrThrow(http.NewRequest("POST", i.testEnvironment.BaseUrl()+"/v1/decrease-product-quantity-in-cart",
				strings.NewReader(template["body"])))
			accessToken := testhelpers.TestGenerateAccessToken(uuid.New())
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", "Bearer "+accessToken)

			response := utils.GetOrThrow(i.testEnvironment.Client().Do(request))

			body := utils.GetOrThrow(io.ReadAll(response.Body))

			i.Equal(400, response.StatusCode)
			i.JSONEq(fmt.Sprintf(`
				{
					"message": %s
				}
			`, template["error"]), string(body))
		}
	})
}

func TestDecreaseProductQuantityInCartSuite(t *testing.T) {
	suite.Run(t, new(DecreaseProductQuantityInCartSuite))
}
