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

type AddStockSuite struct {
	suite.Suite
	productDAO      daos.ProductDAO
	inventoryDAO    daos.InventoryDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (a *AddStockSuite) SetupSuite() {
	a.testEnvironment = testhelpers.NewTestEnvironment()
	err := a.testEnvironment.Start()
	a.Require().NoError(err)

	a.productDAO = daos.NewProductDAO(a.testEnvironment.PgxPool())
	a.inventoryDAO = daos.NewInventoryDAO(a.testEnvironment.PgxPool())
}

func (a *AddStockSuite) SetupTest() {
	err := a.productDAO.DeletAll()
	a.Require().NoError(err)
}

func (a *AddStockSuite) Test1() {
	a.Run("given that the inventory exists, when adding stock, then it returns 204 and sums stock", func() {
		err := a.productDAO.Create(daos.ProductSchema{
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

		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/admin/add-stock", strings.NewReader(`
			{
				"inventoryId": "cf23ee55-88c0-4898-ada4-15645c75645d",
				"stock": 8
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

		inventorySchema, err := a.inventoryDAO.FindOneByProductId(uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"))
		a.Require().NoError(err)
		a.Require().NotNil(inventorySchema)
		a.Require().True(utils.IsValidUUID(inventorySchema.Id.String()))
		a.Require().Equal("c0981e5b-9cb7-4623-9713-55db0317dc1a", inventorySchema.ProductId.String())
		a.Require().Equal(int32(12), inventorySchema.StockQuantity)
		a.Require().WithinDuration(time.Now(), inventorySchema.CreatedAt, 5*time.Second)
	})
}

func (a *AddStockSuite) Test2() {
	a.Run("given that the inventory does not exists, when adding stock, then it returns 409", func() {
		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/admin/add-stock", strings.NewReader(`
			{
				"inventoryId": "cf23ee55-88c0-4898-ada4-15645c75645d",
				"stock": 8
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
				"message": "inventory not found"
			}
		`, string(body))
	})
}

func (a *AddStockSuite) Test3() {
	a.Run("given that the inventory exists, when adding stock and value is zero, then it returns 409", func() {
		request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/admin/add-stock", strings.NewReader(`
			{
				"inventoryId": "cf23ee55-88c0-4898-ada4-15645c75645d",
				"stock": 0
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
				"message": "stock quantity must be higher than zero"
			}
		`, string(body))
	})
}

func (a *AddStockSuite) Test4() {
	a.Run("when adding product and body is invalid, then returns 400", func() {
		templates := []map[string]string{
			{
				"body": `{}`,
				"error": `[
					"inventoryId is required",
					"stock is required"
				]`,
			},
			{
				"body": `{
					"inventoryId": null,
					"stock": null
				}`,
				"error": `[
					"inventoryId is required",
					"stock is required"
				]`,
			},
			{
				"body": `{
					"inventoryId": "",
					"stock": ""
				}`,
				"error": `[
					"inventoryId must be uuidv4",
					"stock must be integer"
				]`,
			},
			{
				"body": `{
					"inventoryId": " ",
					"stock": " "
				}`,
				"error": `[
					"inventoryId must be uuidv4",
					"stock must be integer"
				]`,
			},
			{
				"body": `{
					"inventoryId": 1,
					"stock": 1
				}`,
				"error": `[
					"inventoryId must be uuidv4"
				]`,
			},
			{
				"body": `{
					"inventoryId": 1.5,
					"stock": 1.5
				}`,
				"error": `[
					"inventoryId must be uuidv4",
					"stock must be integer"
				]`,
			},
			{
				"body": `{
					"inventoryId": -1,
					"stock": -1
				}`,
				"error": `[
					"inventoryId must be uuidv4",
					"stock must be positive"
				]`,
			},
			{
				"body": `{
					"inventoryId": true,
					"stock": false
				}`,
				"error": `[
					"inventoryId must be uuidv4",
					"stock must be integer"
				]`,
			},
			{
				"body": `{
					"inventoryId": {},
					"stock": {}
				}`,
				"error": `[
					"inventoryId must be uuidv4",
					"stock must be integer"
				]`,
			},
			{
				"body": `{
					"inventoryId": [],
					"stock": []
				}`,
				"error": `[
					"inventoryId must be uuidv4",
					"stock must be integer"
				]`,
			},
		}

		for _, template := range templates {
			request, err := http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/admin/add-stock", strings.NewReader(template["body"]))
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

func TestAddStock(t *testing.T) {
	suite.Run(t, new(AddStockSuite))
}
