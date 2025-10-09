package usecases

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUsecaseInput struct {
	Name     string
	Email    string
	Password string
}

type RegisterUsecase struct {
	customerDAO daos.CustomerDAO
}

func NewRegisterUsecase(customerDAO daos.CustomerDAO) RegisterUsecase {
	return RegisterUsecase{customerDAO}
}

func (r RegisterUsecase) Execute(input RegisterUsecaseInput) error {
	if len(input.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}

	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		return errors.New("email address is invalid")
	}

	if len(input.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	customerSchema, err := r.customerDAO.FindOneByEmail(input.Email)
	if err != nil {
		return err
	}

	if customerSchema != nil {
		return errors.New("this email address has already been taken by someone")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return r.customerDAO.Create(daos.CustomerSchema{
		Id:        uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
	})
}
