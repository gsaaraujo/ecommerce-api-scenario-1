package apitests_test

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	testhelpers "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/test_helpers"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type RegisterSuite struct {
	suite.Suite
	customerDAO     daos.CustomerDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (r *RegisterSuite) SetupSuite() {
	r.testEnvironment = testhelpers.NewTestEnvironment()
	err := r.testEnvironment.Start()
	r.Require().NoError(err)

	r.customerDAO = daos.NewCustomerDAO(r.testEnvironment.PgxPool())
}

func (r *RegisterSuite) SetupTest() {
	err := r.customerDAO.DeletAll()
	r.Require().NoError(err)
}

func (r *RegisterSuite) Test1() {
	r.Run("given that the customer is not already registered, when registering, then returns 204", func() {
		response, err := r.testEnvironment.Client().Post(r.testEnvironment.BaseUrl()+"/v1/register", "application/json", strings.NewReader(`
			{
				"name": "John Doe",
				"email": "john.doe@gmail.com",
				"password": "123456"
			}
		`))
		r.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		r.Require().NoError(err)
		r.Equal(204, response.StatusCode)
		r.Equal("", string(body))

		customerSchema, err := r.customerDAO.FindOneByEmail("john.doe@gmail.com")
		r.Require().NoError(err)
		r.Require().NotNil(customerSchema)
		r.Require().True(utils.IsValidUUID(customerSchema.Id.String()))
		r.Require().Equal("John Doe", customerSchema.Name)
		r.Require().Equal("john.doe@gmail.com", customerSchema.Email)
		err = bcrypt.CompareHashAndPassword([]byte(customerSchema.Password), []byte("123456"))
		r.Require().NoError(err)
		r.Require().WithinDuration(time.Now().UTC(), customerSchema.CreatedAt, 5*time.Second)
	})
}

func (r *RegisterSuite) Test2() {
	r.Run("given that the email has already been taken by someone, when registering, then returns 409", func() {
		err := r.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		r.Require().NoError(err)

		response, err := r.testEnvironment.Client().Post(r.testEnvironment.BaseUrl()+"/v1/register", "application/json", strings.NewReader(`
			{
				"name": "John Doe Smith",
				"email": "john.doe@gmail.com",
				"password": "123456"
			}
		`))
		r.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		r.Require().NoError(err)
		r.Equal(409, response.StatusCode)
		r.JSONEq(`
			{
				"message": "this email address has already been taken by someone"
			}
		`, string(body))
	})
}

func (r *RegisterSuite) Test3() {
	r.Run("when registering and name is less than 2 characters, then returns 409", func() {
		response, err := r.testEnvironment.Client().Post(r.testEnvironment.BaseUrl()+"/v1/register", "application/json", strings.NewReader(`
			{
				"name": "J",
				"email": "john.doe@gmail.com",
				"password": "123456"
			}
		`))
		r.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		r.Require().NoError(err)
		r.Equal(409, response.StatusCode)
		r.JSONEq(`
			{
				"message": "name must be at least 2 characters"
			}
		`, string(body))
	})
}

func (r *RegisterSuite) Test4() {
	r.Run("when registering and email is invalid, then returns 409", func() {
		response, err := r.testEnvironment.Client().Post(r.testEnvironment.BaseUrl()+"/v1/register", "application/json", strings.NewReader(`
			{
				"name": "John Doe",
				"email": "john",
				"password": "123456"
			}
		`))
		r.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		r.Require().NoError(err)
		r.Equal(409, response.StatusCode)
		r.JSONEq(`
			{
				"message": "email address is invalid"
			}
		`, string(body))
	})
}

func (r *RegisterSuite) Test5() {
	r.Run("when registering and password is less than 6 characters, then returns 409", func() {
		response, err := r.testEnvironment.Client().Post(r.testEnvironment.BaseUrl()+"/v1/register", "application/json", strings.NewReader(`
			{
				"name": "John Doe",
				"email": "john.doe@gmail.com",
				"password": "123"
			}
		`))
		r.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		r.Require().NoError(err)
		r.Equal(409, response.StatusCode)
		r.JSONEq(`
			{
				"message": "password must be at least 6 characters"
			}
		`, string(body))
	})
}

func (r *RegisterSuite) Test6() {
	r.Run("when registering in and body is invalid, then returns 400", func() {
		templates := []map[string]string{
			{
				"body": `{}`,
				"error": `[
					"name is required",
					"email is required",
					"password is required"
				]`,
			},
			{
				"body": `{
					"name": null,
					"email": null,
					"password": null
				}`,
				"error": `[
					"name is required",
					"email is required",
					"password is required"
				]`,
			},
			{
				"body": `{
					"name": "",
					"email": "",
					"password": ""
				}`,
				"error": `[
					"name must not be empty",
					"email must not be empty",
					"password must not be empty"
				]`,
			},
			{
				"body": `{
					"name": " ",
					"email": " ",
					"password": " "
				}`,
				"error": `[
					"name must not be empty",
					"email must not be empty",
					"password must not be empty"
				]`,
			},
			{
				"body": `{
					"name": 1,
					"email": 1,
					"password": 1
				}`,
				"error": `[
					"name must be string",
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"name": 1.5,
					"email": 1.5,
					"password": 1.5
				}`,
				"error": `[
					"name must be string",
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"name": -1,
					"email": -1,
					"password": -1
				}`,
				"error": `[
					"name must be string",
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"name": true,
					"email": true,
					"password": true
				}`,
				"error": `[
					"name must be string",
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"name": {},
					"email": {},
					"password": {}
				}`,
				"error": `[
					"name must be string",
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"name": [],
					"email": [],
					"password": []
				}`,
				"error": `[
					"name must be string",
					"email must be string",
					"password must be string"
				]`,
			},
		}

		for _, template := range templates {
			response, err := r.testEnvironment.Client().Post(r.testEnvironment.BaseUrl()+"/v1/register", "application/json", strings.NewReader(template["body"]))

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

func TestRegister(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}
