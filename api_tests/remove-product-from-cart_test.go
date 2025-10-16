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

type RemoveProductFromCartSuite struct {
	suite.Suite
	productDAO      daos.ProductDAO
	cartDAO         daos.CartDAO
	cartItemDAO     daos.CartItemDAO
	customerDAO     daos.CustomerDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (r *RemoveProductFromCartSuite) SetupSuite() {
	r.testEnvironment = testhelpers.NewTestEnvironment()
	err := r.testEnvironment.Start()
	r.Require().NoError(err)

	r.productDAO = daos.NewProductDAO(r.testEnvironment.PgxPool())
	r.cartDAO = daos.NewCartDAO(r.testEnvironment.PgxPool())
	r.cartItemDAO = daos.NewCartItemDAO(r.testEnvironment.PgxPool())
	r.customerDAO = daos.NewCustomerDAO(r.testEnvironment.PgxPool())
}

func (r *RemoveProductFromCartSuite) SetupTest() {
	err := r.customerDAO.DeletAll()
	r.Require().NoError(err)

	err = r.productDAO.DeletAll()
	r.Require().NoError(err)

	err = r.cartDAO.DeletAll()
	r.Require().NoError(err)

	err = r.cartItemDAO.DeletAll()
	r.Require().NoError(err)
}

func (r *RemoveProductFromCartSuite) Test1() {
	r.Run("given that the product is in cart, when removing product from cart, then returns 204 and updates cart", func() {
		err := r.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		r.Require().NoError(err)
		err = r.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		r.Require().NoError(err)
		err = r.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		r.Require().NoError(err)
		err = r.cartItemDAO.Create(daos.CartItemSchema{
			Id:        uuid.MustParse("b999870f-f969-4d24-8955-499dbf3c689e"),
			CartId:    uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			ProductId: uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Quantity:  8,
			CreatedAt: time.Now().UTC(),
		})
		r.Require().NoError(err)

		request, err := http.NewRequest("POST", r.testEnvironment.BaseUrl()+"/v1/remove-product-from-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a"
			}
		`))
		r.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		r.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := r.testEnvironment.Client().Do(request)
		r.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		r.Require().NoError(err)
		r.Equal(204, response.StatusCode)
		r.Equal("", string(body))

		cartSchema, err := r.cartDAO.FindOneByCustomerId(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		r.Require().NoError(err)
		r.Require().NotNil(cartSchema)

		cartItemSchema, err := r.cartItemDAO.FindAllByCartId(cartSchema.Id)
		r.Require().NoError(err)
		r.Require().Equal(0, len(cartItemSchema))
	})
}

func (r *RemoveProductFromCartSuite) Test2() {
	r.Run("given that the product is not found in cart, when removing product from cart, then returns 409", func() {
		err := r.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		r.Require().NoError(err)
		err = r.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		r.Require().NoError(err)
		err = r.cartDAO.Create(daos.CartSchema{
			Id:         uuid.MustParse("bb8357b2-b978-4675-9521-ef2da0bd1747"),
			CustomerId: uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			CreatedAt:  time.Now().UTC(),
		})
		r.Require().NoError(err)

		request, err := http.NewRequest("POST", r.testEnvironment.BaseUrl()+"/v1/remove-product-from-cart", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a"
			}
		`))
		r.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		r.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := r.testEnvironment.Client().Do(request)
		r.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		r.Require().NoError(err)
		r.Equal(409, response.StatusCode)
		r.JSONEq(`
			{
				"message": "product not found in cart"
			}
		`, string(body))
	})
}

func (r *RemoveProductFromCartSuite) Test3() {
	r.Run("when removing product from cart and body is invalid, then returns 400", func() {
		templates := []map[string]string{
			{
				"body": `{}`,
				"error": `[
					"productId is required"
				]`,
			},
			{
				"body": `{
					"productId": null
				}`,
				"error": `[
					"productId is required"
				]`,
			},
			{
				"body": `{
					"productId": ""
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": " "
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": 1
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": 1.5
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": -1
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": true
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": {}
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"productId": []
				}`,
				"error": `[
					"productId must be uuidv4"
				]`,
			},
		}

		for _, template := range templates {
			request, err := http.NewRequest("POST", r.testEnvironment.BaseUrl()+"/v1/remove-product-from-cart", strings.NewReader(template["body"]))
			r.Require().NoError(err)
			accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
			r.Require().NoError(err)
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", "Bearer "+accessToken)

			response, err := r.testEnvironment.Client().Do(request)
			r.Require().NoError(err)

			body, err := io.ReadAll(response.Body)
			r.Require().NoError(err)

			r.Equal(400, response.StatusCode)
			r.JSONEq(fmt.Sprintf(`
				{
					"message": %s
				}
			`, template["error"]), string(body))
		}
	})
}

func TestRemoveProductFromCart(t *testing.T) {
	suite.Run(t, new(RemoveProductFromCartSuite))
}
