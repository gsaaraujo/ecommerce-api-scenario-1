package apitests_test

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	testhelpers "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/test_helpers"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/stretchr/testify/suite"
)

type LoginSuite struct {
	suite.Suite
	customerDAO     daos.CustomerDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (l *LoginSuite) SetupSuite() {
	l.testEnvironment = testhelpers.NewTestEnvironment()
	err := l.testEnvironment.Start()
	l.Require().NoError(err)

	l.customerDAO = daos.NewCustomerDAO(l.testEnvironment.PgxPool())
}

func (l *LoginSuite) SetupTest() {
	err := l.customerDAO.DeletAll()
	l.Require().NoError(err)
}

func (l *LoginSuite) Test1() {
	l.Run("given that the customer is already registered, when logging in, then returns 200", func() {
		err := l.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		l.Require().NoError(err)

		response, err := l.testEnvironment.Client().Post(l.testEnvironment.BaseUrl()+"/v1/login", "application/json", strings.NewReader(`
			{
				"email": "john.doe@gmail.com",
				"password": "123456"
			}
		`))
		l.Require().NoError(err)

		l.Equal(200, response.StatusCode)
		responseBody, err := utils.ParseJSONBody[map[string]map[string]any](response.Body)
		l.Require().NoError(err)
		customerId := responseBody["data"]["customerId"].(string)
		accessToken := responseBody["data"]["accessToken"].(string)
		l.True(utils.IsValidUUID(customerId))
		l.NotEqual("", accessToken)

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
			return []byte("81c4a8d5b2554de4ba736e93255ba633"), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		l.Require().NoError(err)

		subject, err := token.Claims.GetSubject()
		l.Require().NoError(err)
		issuedAt, err := token.Claims.GetIssuedAt()
		l.Require().NoError(err)
		expirationTime, err := token.Claims.GetExpirationTime()
		l.Require().NoError(err)

		l.Equal("f59207c8-e837-4159-b67d-78c716510747", subject)
		l.WithinDuration(time.Now().UTC(), issuedAt.Time, 5*time.Second)
		l.WithinDuration(time.Now().UTC().Add(30*time.Minute), expirationTime.Time, 5*time.Second)
	})
}

func (l *LoginSuite) Test2() {
	l.Run("given that the customer is already registered, when logging in and password is incorrect, then returns 409", func() {
		err := l.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		l.Require().NoError(err)

		response, err := l.testEnvironment.Client().Post(l.testEnvironment.BaseUrl()+"/v1/login", "application/json", strings.NewReader(`
			{
				"email": "john.doe@gmail.com",
				"password": "abc123"
			}
		`))
		l.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		l.Require().NoError(err)

		l.Equal(409, response.StatusCode)
		l.JSONEq(`
			{
				"message": "email or password is incorrect"
			}
		`, string(body))
	})
}

func (l *LoginSuite) Test3() {
	l.Run("given that the customer is not registered, when logging in, then returns 409", func() {
		response, err := l.testEnvironment.Client().Post(l.testEnvironment.BaseUrl()+"/v1/login", "application/json", strings.NewReader(`
			{
				"email": "john.doe@gmail.com",
				"password": "123456"
			}
		`))
		l.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		l.Require().NoError(err)

		l.Equal(409, response.StatusCode)
		l.JSONEq(`
			{
				"message": "email or password is incorrect"
			}
		`, string(body))
	})
}

func (l *LoginSuite) Test4() {
	l.Run("when logging in and email address is invalid, then returns 409", func() {
		response, err := l.testEnvironment.Client().Post(l.testEnvironment.BaseUrl()+"/v1/login", "application/json", strings.NewReader(`
			{
				"email": "john",
				"password": "123456"
			}
		`))
		l.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		l.Require().NoError(err)

		l.Equal(409, response.StatusCode)
		l.JSONEq(`
			{
				"message": "email address is invalid"
			}
		`, string(body))
	})
}

func (l *LoginSuite) Test5() {
	l.Run("when adding logging in and body is invalid, then returns 400", func() {
		templates := []map[string]string{
			{
				"body": `{}`,
				"error": `[
					"email is required",
					"password is required"
				]`,
			},
			{
				"body": `{
					"email": null,
					"password": null
				}`,
				"error": `[
					"email is required",
					"password is required"
				]`,
			},
			{
				"body": `{
					"email": "",
					"password": ""
				}`,
				"error": `[
					"email must not be empty",
					"password must not be empty"
				]`,
			},
			{
				"body": `{
					"email": " ",
					"password": " "
				}`,
				"error": `[
					"email must not be empty",
					"password must not be empty"
				]`,
			},
			{
				"body": `{
					"email": 1,
					"password": 1
				}`,
				"error": `[
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"email": 1.5,
					"password": 1.5
				}`,
				"error": `[
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"email": -1,
					"password": -1
				}`,
				"error": `[
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"email": true,
					"password": true
				}`,
				"error": `[
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"email": {},
					"password": {}
				}`,
				"error": `[
					"email must be string",
					"password must be string"
				]`,
			},
			{
				"body": `{
					"email": [],
					"password": []
				}`,
				"error": `[
					"email must be string",
					"password must be string"
				]`,
			},
		}

		for _, template := range templates {
			response, err := l.testEnvironment.Client().Post(l.testEnvironment.BaseUrl()+"/v1/login", "application/json", strings.NewReader(template["body"]))

			l.Require().NoError(err)

			body, err := io.ReadAll(response.Body)
			l.Require().NoError(err)
			l.Equal(400, response.StatusCode)
			l.JSONEq(fmt.Sprintf(`
				{
					"message": %s
				}
			`, template["error"]), string(body))
		}
	})
}

func TestLogin(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}
