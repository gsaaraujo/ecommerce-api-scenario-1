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

type PublishProductSuite struct {
	suite.Suite
	productDAO      daos.ProductDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (p *PublishProductSuite) SetupSuite() {
	p.testEnvironment = testhelpers.NewTestEnvironment()
	err := p.testEnvironment.Start()
	p.Require().NoError(err)

	p.productDAO = daos.NewProductDAO(p.testEnvironment.PgxPool())
}

func (p *PublishProductSuite) SetupTest() {
	err := p.productDAO.DeletAll()
	p.Require().NoError(err)
}

func (p *PublishProductSuite) Test1() {
	p.Run("given that the product exists, when publishing, then returns 204 and changes product status to 'published'", func() {
		err := p.productDAO.Create(daos.ProductSchema{
			Id:          uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"),
			Name:        "ErgoClick Pro Wireless Mouse",
			Description: utils.NewPointer("Ergonomically designed wireless optical mouse ..."),
			Price:       2999,
			CreatedAt:   time.Now().UTC(),
		})
		p.Require().NoError(err)
		request, err := http.NewRequest("POST", p.testEnvironment.BaseUrl()+"/v1/admin/publish-product", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a"
			}
		`))
		p.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
		p.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := p.testEnvironment.Client().Do(request)
		p.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		p.Require().NoError(err)
		p.Equal(204, response.StatusCode)
		p.Equal("", string(body))

		productSchema, err := p.productDAO.FindOneById(uuid.MustParse("c0981e5b-9cb7-4623-9713-55db0317dc1a"))
		p.Require().NoError(err)
		p.Require().NotNil(productSchema)
		p.Require().Equal("published", productSchema.Status)
	})
}

func (p *PublishProductSuite) Test2() {
	p.Run("given that the product does not exist, when publishing, then returns 409", func() {
		request, err := http.NewRequest("POST", p.testEnvironment.BaseUrl()+"/v1/admin/publish-product", strings.NewReader(`
			{
				"productId": "c0981e5b-9cb7-4623-9713-55db0317dc1a"
			}
		`))
		p.Require().NoError(err)
		accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
		p.Require().NoError(err)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response, err := p.testEnvironment.Client().Do(request)
		p.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		p.Require().NoError(err)
		p.Equal(409, response.StatusCode)
		p.JSONEq(`
			{
				"message": "product not found"
			}
		`, string(body))
	})
}

func (p *PublishProductSuite) Test4() {
	p.Run("when publishing product and body is invalid, then returns 400", func() {
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
			request, err := http.NewRequest("POST", p.testEnvironment.BaseUrl()+"/v1/admin/publish-product", strings.NewReader(template["body"]))
			p.Require().NoError(err)
			accessToken, err := testhelpers.TestGenerateAccessToken(uuid.New())
			p.Require().NoError(err)
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", "Bearer "+accessToken)

			response, err := p.testEnvironment.Client().Do(request)
			p.Require().NoError(err)

			body, err := io.ReadAll(response.Body)
			p.Require().NoError(err)

			p.Equal(400, response.StatusCode)
			p.JSONEq(fmt.Sprintf(`
				{
					"message": %s
				}
			`, template["error"]), string(body))
		}
	})
}

func TestPublishProduct(t *testing.T) {
	suite.Run(t, new(PublishProductSuite))
}
