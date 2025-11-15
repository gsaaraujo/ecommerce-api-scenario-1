package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/gateways"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/redis/go-redis/v9"
)

type AddAddressUsecaseInput struct {
	CustomerId   uuid.UUID
	City         string
	State        string
	ZipCode      string
	StreetName   string
	StreetNumber string
}

type AddAddressUsecase struct {
	redisClient        *redis.Client
	addressDAO         daos.AddressDAO
	httpZipCodeGateway gateways.HttpZipCodeGateway
}

func NewAddAddressUsecase(redisClient *redis.Client, addressDAO daos.AddressDAO, httpZipCodeGateway gateways.HttpZipCodeGateway) AddAddressUsecase {
	return AddAddressUsecase{redisClient, addressDAO, httpZipCodeGateway}
}

func (a *AddAddressUsecase) Execute(input AddAddressUsecaseInput) error {
	validStates := []string{"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA", "HI", "ID", "IL", "IN", "IA", "KS", "KY",
		"LA", "ME", "MD", "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA",
		"RI", "SC", "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY", "DC"}

	if !slices.Contains(validStates, strings.ToUpper(input.State)) {
		return errors.New("state must be a valid 2-letter U.S. abbreviation (e.g. NY, CA)")
	}

	zipRegex := regexp.MustCompile(`^[0-9]{5}$`)

	if !zipRegex.MatchString(input.ZipCode) {
		return errors.New("ZIP code is invalid. It must be 5 digits (e.g. 12345)")
	}

	if _, err := strconv.ParseUint(input.StreetNumber, 10, 64); err != nil {
		return errors.New("street number must contain only digits (0-9)")
	}

	type Location struct {
		City  string `json:"city"`
		State string `json:"state"`
	}
	var location *Location = nil

	cachedZipCode, err := a.redisClient.Get(context.Background(), "zip_codes:"+input.ZipCode).Result()
	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	if cachedZipCode == "" {
		httpZipCodeResponse := utils.GetOrThrow(a.httpZipCodeGateway.Get(input.ZipCode))

		if httpZipCodeResponse != nil {
			location = &Location{
				City:  httpZipCodeResponse.City,
				State: httpZipCodeResponse.State,
			}

			locationJson := utils.GetOrThrow(json.Marshal(location))
			utils.ThrowOnError(a.redisClient.Set(context.Background(), "zip_codes:"+input.ZipCode, string(locationJson), 0).Err())
		}
	} else {
		var locationJson Location
		utils.ThrowOnError(json.Unmarshal([]byte(cachedZipCode), &locationJson))

		location = &Location{
			City:  locationJson.City,
			State: locationJson.State,
		}
	}

	if location == nil {
		return errors.New("ZIP code does not match any location")
	}

	if location.City != input.City || location.State != input.State {
		return errors.New("ZIP code location does not match with provided city and state")
	}

	addressLine := fmt.Sprintf("%s %s, %s, %s %s", input.StreetNumber, input.StreetName, input.City, input.State, input.ZipCode)
	isThereDefaultAddress := a.addressDAO.FindOneByIsDefault(true)

	a.addressDAO.Create(daos.AddressSchema{
		Id:          uuid.New(),
		CustomerId:  input.CustomerId,
		IsDefault:   !isThereDefaultAddress,
		Street:      input.StreetName,
		Number:      input.StreetNumber,
		City:        input.City,
		State:       input.State,
		ZipCode:     input.ZipCode,
		AddressLine: addressLine,
		CreatedAt:   time.Now().UTC(),
	})

	return nil
}
