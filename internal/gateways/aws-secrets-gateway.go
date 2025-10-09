package gateways

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSecretsGateway struct {
	secretsClient *secretsmanager.Client
}

func NewAwsSecretsGateway(secretsClient *secretsmanager.Client) AwsSecretsGateway {
	return AwsSecretsGateway{
		secretsClient: secretsClient,
	}
}

func (a *AwsSecretsGateway) Get(key string) (string, error) {
	if _, ok := os.LookupEnv("AWS_SECRET_MANAGER_NAME"); !ok {
		return "", errors.New("AWS_SECRET_MANAGER_NAME environment variable not found")
	}

	secretValue, err := a.secretsClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("AWS_SECRET_MANAGER_NAME")),
	})
	if err != nil {
		return "", err
	}

	var secret map[string]any
	err = json.Unmarshal([]byte(*secretValue.SecretString), &secret)
	if err != nil {
		return "", err
	}

	value, exists := secret[key]

	if !exists {
		return "", fmt.Errorf("%s secret not found", key)
	}

	return fmt.Sprint(value), nil
}
