package gateways

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
)

type HttpZipCodeSuccessResponse struct {
	ZipCode string `json:"zip_code"`
	City    string `json:"city"`
	State   string `json:"state"`
}

type HttpZipCodeFailResponse struct {
	ErrorMsg string `json:"error_msg"`
}

type HttpZipCodeGateway struct {
	awsSecretsGateway AwsSecretsGateway
}

func NewHttpZipCodeGateway(awsSecretsGateway AwsSecretsGateway) HttpZipCodeGateway {
	return HttpZipCodeGateway{awsSecretsGateway}
}

func (h *HttpZipCodeGateway) Get(zipCode string) (*HttpZipCodeSuccessResponse, error) {
	if _, ok := os.LookupEnv("ZIPCODE_URL"); !ok {
		return nil, errors.New("ZIPCODE_URL environment variable not found")
	}

	zipCodeToken, err := h.awsSecretsGateway.Get("ZIPCODE_TOKEN")
	if err != nil {
		return nil, err
	}

	response, err := http.Get(fmt.Sprintf("%s/rest/%s/info.json/%s/degrees", os.Getenv("ZIPCODE_URL"), zipCodeToken, zipCode))
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 200 {
		responseBody, err := utils.ParseJSONBody[HttpZipCodeSuccessResponse](response.Body)
		if err != nil {
			return nil, err
		}

		return &responseBody, nil
	}

	if response.StatusCode == 404 {
		responseBody, err := utils.ParseJSONBody[HttpZipCodeFailResponse](response.Body)
		if err != nil {
			return nil, err
		}

		if responseBody.ErrorMsg == fmt.Sprintf(`Zip code "%s" not found.`, zipCode) {
			return nil, nil
		}
	}

	return nil, fmt.Errorf("httpZipCodeError, HTTP status %d", response.StatusCode)
}
