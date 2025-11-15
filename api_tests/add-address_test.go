package apitests_test

import (
	"context"
	"encoding/json"
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

type AddAddressSuite struct {
	suite.Suite
	customerDAO     daos.CustomerDAO
	addressDAO      daos.AddressDAO
	testEnvironment *testhelpers.TestEnvironment
}

func (a *AddAddressSuite) SetupSuite() {
	a.testEnvironment = testhelpers.NewTestEnvironment()
	a.testEnvironment.Start()

	a.customerDAO = daos.NewCustomerDAO(a.testEnvironment.PgxPool())
	a.addressDAO = daos.NewAddressDAO(a.testEnvironment.PgxPool())
}

func (a *AddAddressSuite) SetupTest() {
	a.customerDAO.DeletAll()
	a.addressDAO.DeletAll()
	utils.ThrowOnError(a.testEnvironment.RedisClient().FlushDB(context.Background()).Err())
}

func (a *AddAddressSuite) Test1() {
	a.Run("given customer has no address, when adding, returns 201 and a new address is created and it's the default", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/73301/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "73301",
						"city": "Austin",
						"state": "TX"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "73301",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		addressesSchema := a.addressDAO.FindAllByCustomerId(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NotEmpty(addressesSchema)
		a.Require().True(utils.IsValidUUID(addressesSchema[0].Id.String()))
		a.Require().Equal("f59207c8-e837-4159-b67d-78c716510747", addressesSchema[0].CustomerId.String())
		a.Require().Equal(true, addressesSchema[0].IsDefault)
		a.Require().Equal("Delivery Road", addressesSchema[0].Street)
		a.Require().Equal("Austin", addressesSchema[0].City)
		a.Require().Equal("TX", addressesSchema[0].State)
		a.Require().Equal("321", addressesSchema[0].Number)
		a.Require().Equal("73301", addressesSchema[0].ZipCode)
		a.Require().Equal("321 Delivery Road, Austin, TX 73301", addressesSchema[0].AddressLine)
	})
}

func (a *AddAddressSuite) Test2() {
	a.Run("given customer already has a default address, when adding, returns 201 and a new address is created and it's not default", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		a.addressDAO.Create(daos.AddressSchema{
			Id:          uuid.New(),
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
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/73301/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "73301",
						"city": "Austin",
						"state": "TX"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "73301",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		addressesSchema := a.addressDAO.FindAllByCustomerId(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NotEmpty(addressesSchema)
		a.Require().True(utils.IsValidUUID(addressesSchema[1].Id.String()))
		a.Require().Equal("f59207c8-e837-4159-b67d-78c716510747", addressesSchema[1].CustomerId.String())
		a.Require().Equal(false, addressesSchema[1].IsDefault)
		a.Require().Equal("Delivery Road", addressesSchema[1].Street)
		a.Require().Equal("Austin", addressesSchema[1].City)
		a.Require().Equal("TX", addressesSchema[1].State)
		a.Require().Equal("321", addressesSchema[1].Number)
		a.Require().Equal("73301", addressesSchema[1].ZipCode)
		a.Require().Equal("321 Delivery Road, Austin, TX 73301", addressesSchema[1].AddressLine)
	})
}

func (a *AddAddressSuite) Test3() {
	a.Run("given customer already has a address and it's not default, when adding, returns 201 and a new address is created as default", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		a.addressDAO.Create(daos.AddressSchema{
			Id:          uuid.New(),
			CustomerId:  uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			IsDefault:   false,
			Street:      "Maple Grove Lane",
			Number:      "4767",
			City:        "Austin",
			State:       "TX",
			ZipCode:     "78739",
			AddressLine: "4767 Maple Grove Lane, Austin, TX 78739",
			CreatedAt:   time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/73301/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "73301",
						"city": "Austin",
						"state": "TX"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "73301",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		addressesSchema := a.addressDAO.FindAllByCustomerId(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		a.Require().NotEmpty(addressesSchema)
		a.Require().True(utils.IsValidUUID(addressesSchema[1].Id.String()))
		a.Require().Equal("f59207c8-e837-4159-b67d-78c716510747", addressesSchema[1].CustomerId.String())
		a.Require().Equal(true, addressesSchema[1].IsDefault)
		a.Require().Equal("Delivery Road", addressesSchema[1].Street)
		a.Require().Equal("Austin", addressesSchema[1].City)
		a.Require().Equal("TX", addressesSchema[1].State)
		a.Require().Equal("321", addressesSchema[1].Number)
		a.Require().Equal("73301", addressesSchema[1].ZipCode)
		a.Require().Equal("321 Delivery Road, Austin, TX 73301", addressesSchema[1].AddressLine)
	})
}

func (a *AddAddressSuite) Test4() {
	a.Run("given that the zip code info from external api is not cached, when adding, returns 201 and cache it", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/73301/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "73301",
						"city": "Austin",
						"state": "TX"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "73301",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		zipCode := utils.GetOrThrow(a.testEnvironment.RedisClient().Get(context.Background(), "zip_codes:73301").Result())
		a.Require().JSONEq(`
			{
				"city": "Austin",
				"state": "TX"
			}
		`, zipCode)
	})
}

func (a *AddAddressSuite) Test5() {
	a.Run("given that the zip code info from external api is already cached, when adding, returns 201", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		location := utils.GetOrThrow(json.Marshal(map[string]any{
			"city":  "Austin",
			"state": "TX",
		}))
		utils.ThrowOnError(a.testEnvironment.RedisClient().Set(context.Background(), "zip_codes:73301", location, 0).Err())

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "73301",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(204, response.StatusCode)
		a.Equal("", string(body))

		zipCode := utils.GetOrThrow(a.testEnvironment.RedisClient().Get(context.Background(), "zip_codes:73301").Result())
		a.Require().JSONEq(`
			{
				"city": "Austin",
				"state": "TX"
			}
		`, zipCode)
	})
}

func (a *AddAddressSuite) Test6() {
	a.Run("when adding and ZIP code does not match any location, returns 409", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/11111/degrees"
				},
				"response": {
					"status": 404,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"error_code": 404,
						"error_msg": "Zip code \"11111\" not found."
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "11111",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "ZIP code does not match any location"
			}
		`, string(body))
	})
}

func (a *AddAddressSuite) Test7() {
	a.Run("when adding and ZIP code does not match any location, returns 409", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/90001/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "90001",
						"city": "Los Angeles",
						"state": "LA"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "90001",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "ZIP code location does not match with provided city and state"
			}
		`, string(body))
	})
}

func (a *AddAddressSuite) Test8() {
	a.Run("when adding and state is invalid, returns 409", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/73301/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "73301",
						"city": "Austin",
						"state": "TX"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "AB",
				"zipCode": "73301",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "state must be a valid 2-letter U.S. abbreviation (e.g. NY, CA)"
			}
		`, string(body))
	})
}

func (a *AddAddressSuite) Test9() {
	a.Run("when adding and ZIP code is invalid, returns 409", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/73301/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "73301",
						"city": "Austin",
						"state": "TX"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "12345678",
				"streetName": "Delivery Road",
				"streetNumber": "321"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "ZIP code is invalid. It must be 5 digits (e.g. 12345)"
			}
		`, string(body))
	})
}

func (a *AddAddressSuite) Test10() {
	a.Run("when adding and ZIP code is invalid, returns 409", func() {
		a.customerDAO.Create(daos.CustomerSchema{
			Id:        uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"),
			Name:      "John Doe",
			Email:     "john.doe@gmail.com",
			Password:  "$2a$10$asLIHej6kxd3Fsdc76QHieBugwCGvsYJeLiZmP1K7/t1GbIbUy.pK",
			CreatedAt: time.Now().UTC(),
		})
		mockRes := utils.GetOrThrow(http.Post(a.testEnvironment.WiremockContainerUrl()+"/__admin/mappings", "application/json", strings.NewReader(`
			{
				"request": {
					"method": "GET",
					"url": "/rest/a7416146283d464294cebea38d5cb5ff/info.json/73301/degrees"
				},
				"response": {
					"status": 200,
					"headers": {
						"Content-Type": "application/json"
					},
					"jsonBody": {
						"zip_code": "73301",
						"city": "Austin",
						"state": "TX"
					}
				}
			}
		`)))
		a.Require().Equal(201, mockRes.StatusCode)

		request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(`
			{
				"city": "Austin",
				"state": "TX",
				"zipCode": "73301",
				"streetName": "Delivery Road",
				"streetNumber": "AAA"
			}
		`)))
		accessToken := testhelpers.TestGenerateAccessToken(uuid.MustParse("f59207c8-e837-4159-b67d-78c716510747"))
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+accessToken)

		response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

		body := utils.GetOrThrow(io.ReadAll(response.Body))
		a.Equal(409, response.StatusCode)
		a.JSONEq(`
			{
				"message": "street number must contain only digits (0-9)"
			}
		`, string(body))
	})
}

func (a *AddAddressSuite) Test11() {
	a.Run("when adding address and body is invalid, then returns 400", func() {
		templates := []map[string]string{
			{
				"body": `{}`,
				"error": `[
					"city is required",
					"state is required",
					"zipCode is required",
					"streetName is required",
					"streetNumber is required"
				]`,
			},
			{
				"body": `{
					"city": null,
					"state": null,
					"zipCode": null,
					"streetName": null,
					"streetNumber": null
				}`,
				"error": `[
					"city is required",
					"state is required",
					"zipCode is required",
					"streetName is required",
					"streetNumber is required"
				]`,
			},
			{
				"body": `{
					"city": "",
					"state": "",
					"zipCode": "",
					"streetName": "",
					"streetNumber": ""
				}`,
				"error": `[
					"city must not be empty",
					"state must not be empty",
					"zipCode must not be empty",
					"streetName must not be empty",
					"streetNumber must not be empty"
				]`,
			},
			{
				"body": `{
					"city": " ",
					"state": " ",
					"zipCode": " ",
					"streetName": " ",
					"streetNumber": " "
				}`,
				"error": `[
					"city must not be empty",
					"state must not be empty",
					"zipCode must not be empty",
					"streetName must not be empty",
					"streetNumber must not be empty"
				]`,
			},
			{
				"body": `{
					"city": 1,
					"state": 1,
					"zipCode": 1,
					"streetName": 1,
					"streetNumber": 1
				}`,
				"error": `[
					"city must be string",
					"state must be string",
					"zipCode must be string",
					"streetName must be string",
					"streetNumber must be string"
				]`,
			},
			{
				"body": `{
					"city": 1.5,
					"state": 1.5,
					"zipCode": 1.5,
					"streetName": 1.5,
					"streetNumber": 1.5
				}`,
				"error": `[
					"city must be string",
					"state must be string",
					"zipCode must be string",
					"streetName must be string",
					"streetNumber must be string"
				]`,
			},
			{
				"body": `{
					"city": -1,
					"state": -1,
					"zipCode": -1,
					"streetName": -1,
					"streetNumber": -1
				}`,
				"error": `[
					"city must be string",
					"state must be string",
					"zipCode must be string",
					"streetName must be string",
					"streetNumber must be string"
				]`,
			},
			{
				"body": `{
					"city": true,
					"state": true,
					"zipCode": true,
					"streetName": false,
					"streetNumber": false
				}`,
				"error": `[
					"city must be string",
					"state must be string",
					"zipCode must be string",
					"streetName must be string",
					"streetNumber must be string"
				]`,
			},
			{
				"body": `{
					"city": {},
					"state": {},
					"zipCode": {},
					"streetName": {},
					"streetNumber": {}
				}`,
				"error": `[
					"city must be string",
					"state must be string",
					"zipCode must be string",
					"streetName must be string",
					"streetNumber must be string"
				]`,
			},
			{
				"body": `{
					"city": [],
					"state": [],
					"zipCode": [],
					"streetName": [],
					"streetNumber": []
				}`,
				"error": `[
					"city must be string",
					"state must be string",
					"zipCode must be string",
					"streetName must be string",
					"streetNumber must be string"
				]`,
			},
		}

		for _, template := range templates {
			request := utils.GetOrThrow(http.NewRequest("POST", a.testEnvironment.BaseUrl()+"/v1/add-address", strings.NewReader(template["body"])))
			accessToken := testhelpers.TestGenerateAccessToken(uuid.New())
			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", "Bearer "+accessToken)

			response := utils.GetOrThrow(a.testEnvironment.Client().Do(request))

			body := utils.GetOrThrow(io.ReadAll(response.Body))

			a.Equal(400, response.StatusCode)
			a.JSONEq(fmt.Sprintf(`
				{
					"message": %s
				}
			`, template["error"]), string(body))
		}
	})
}

func TestAddAddress(t *testing.T) {
	suite.Run(t, new(AddAddressSuite))
}
