package usecases

import (
	"errors"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/gateways"
	"golang.org/x/crypto/bcrypt"
)

type LoginUsecaseInput struct {
	Email    string
	Password string
}

type LoginUsecaseOutput struct {
	CustomerId  uuid.UUID
	AccessToken string
}

type LoginUsecase struct {
	customerDAO       daos.CustomerDAO
	awsSecretsGateway gateways.AwsSecretsGateway
}

func NewLoginUsecase(customerDAO daos.CustomerDAO, awsSecretsGateway gateways.AwsSecretsGateway) LoginUsecase {
	return LoginUsecase{customerDAO, awsSecretsGateway}
}

func (l *LoginUsecase) Execute(input LoginUsecaseInput) (LoginUsecaseOutput, error) {
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		return LoginUsecaseOutput{}, errors.New("email address is invalid")
	}

	customerSchema, err := l.customerDAO.FindOneByEmail(input.Email)
	if err != nil {
		return LoginUsecaseOutput{}, err
	}

	if customerSchema == nil {
		return LoginUsecaseOutput{}, errors.New("email or password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(customerSchema.Password), []byte(input.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return LoginUsecaseOutput{}, errors.New("email or password is incorrect")
	}
	if err != nil {
		return LoginUsecaseOutput{}, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   customerSchema.Id.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(30 * time.Minute)),
	})

	accessTokenSigningKey, err := l.awsSecretsGateway.Get("ACCESS_TOKEN_SIGNING_KEY")
	if err != nil {
		return LoginUsecaseOutput{}, err
	}

	acessTokenSigned, err := accessToken.SignedString([]byte(accessTokenSigningKey))
	if err != nil {
		return LoginUsecaseOutput{}, err
	}

	return LoginUsecaseOutput{
		CustomerId:  customerSchema.Id,
		AccessToken: acessTokenSigned,
	}, nil
}
