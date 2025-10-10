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

type AddProductSuite struct {
	suite.Suite
	productDAO      daos.ProductDAO
	inventoryDAO    daos.InventoryDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (a *AddProductSuite) SetupSuite() {
	a.testEnvironment = testhelpers.NewTestEnvironment()
	err := a.testEnvironment.Start()
	a.Require().NoError(err)

	a.productDAO = daos.NewProductDAO(a.testEnvironment.PgxPool())
	a.inventoryDAO = daos.NewInventoryDAO(a.testEnvironment.PgxPool())
}

func (a *AddProductSuite) SetupTest() {
	err := a.inventoryDAO.DeletAll()
	a.Require().NoError(err)

	err = a.productDAO.DeletAll()
	a.Require().NoError(err)
}

func (a *AddProductSuite) Test1() {
	a.Run("when adding product, then returns 204 and a new product and inventory are created and inventory has stock equals zero", func() {
		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/admin/add-product", strings.NewReader(`
			{
				"name": "ErgoClick Pro Wireless Mouse",
				"description": "Ergonomically designed wireless optical mouse ...",
				"price": 2999
			}
		`))
		a.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
		a.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := a.testEnvironment.Client().Do(request)
		a.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		a.Require().NoError(err)
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		productSchema, err := a.productDAO.FindOneByName("ErgoClick Pro Wireless Mouse")
		a.Require().NoError(err)
		a.Require().NotNil(productSchema)
		a.Require().True(utils.IsValidUUID(productSchema.Id.String()))
		a.Require().Equal("ErgoClick Pro Wireless Mouse", productSchema.Name)
		a.Require().Equal("Ergonomically designed wireless optical mouse ...", *productSchema.Description)
		a.Require().Equal(int64(2999), productSchema.Price)
		a.Require().WithinDuration(time.Now(), productSchema.CreatedAt, 5*time.Second)

		inventorySchema, err := a.inventoryDAO.FindOneByProductId(productSchema.Id)
		a.Require().NoError(err)
		a.Require().NotNil(inventorySchema)
		a.Require().True(utils.IsValidUUID(inventorySchema.Id.String()))
		a.Require().Equal(productSchema.Id, inventorySchema.ProductId)
		a.Require().Equal(int32(0), inventorySchema.StockQuantity)
		a.Require().WithinDuration(time.Now(), inventorySchema.CreatedAt, 5*time.Second)
	})
}

func (a *AddProductSuite) Test2() {
	a.Run("when adding product and the price is zero, then returns 409", func() {
		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/admin/add-product", strings.NewReader(`
			{
				"name": "ErgoClick Pro Wireless Mouse",
				"description": "Ergonomically designed wireless optical mouse ...",
				"price": 0
			}
		`))
		a.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
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
				"message": "the product price cannot be zero"
			}
		`, string(body))
	})
}

func (a *AddProductSuite) Test3() {
	a.Run("when adding product and body is invalid, then returns 400", func() {
		templates := []map[string]string{
			{
				"body": `{}`,
				"error": `[
					"name is required",
					"price is required"
				]`,
			},
			{
				"body": `{
					"name": null,
					"description": null,
					"price": null
				}`,
				"error": `[
					"name is required",
					"price is required"
				]`,
			},
			{
				"body": `{
					"name": "",
					"description": "",
					"price": ""
				}`,
				"error": `[
					"name must not be empty",
					"description must not be empty",
					"price must be integer"
				]`,
			},
			{
				"body": `{
					"name": " ",
					"description": " ",
					"price": " "
				}`,
				"error": `[
					"name must not be empty",
					"description must not be empty",
					"price must be integer"
				]`,
			},
			{
				"body": `{
					"name": 1,
					"description": 1,
					"price": 1
				}`,
				"error": `[
					"name must be string",
					"description must be string"
				]`,
			},
			{
				"body": `{
					"name": 1.5,
					"description": 1.5,
					"price": 1.5
				}`,
				"error": `[
					"name must be string",
					"description must be string",
					"price must be integer"
				]`,
			},
			{
				"body": `{
					"name": -1,
					"description": -1,
					"price": -1
				}`,
				"error": `[
					"name must be string",
					"description must be string",
					"price must be positive"
				]`,
			},
			{
				"body": `{
					"name": true,
					"description": true,
					"price": false
				}`,
				"error": `[
					"name must be string",
					"description must be string",
					"price must be integer"
				]`,
			},
			{
				"body": `{
					"name": {},
					"description": {},
					"price": {}
				}`,
				"error": `[
					"name must be string",
					"description must be string",
					"price must be integer"
				]`,
			},
			{
				"body": `{
					"name": [],
					"description": [],
					"price": []
				}`,
				"error": `[
					"name must be string",
					"description must be string",
					"price must be integer"
				]`,
			},
		}

		for _, template := range templates {
			request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/admin/add-product", strings.NewReader(template["body"]))
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

func TestAddProduct(t *testing.T) {
	suite.Run(t, new(AddProductSuite))
}
