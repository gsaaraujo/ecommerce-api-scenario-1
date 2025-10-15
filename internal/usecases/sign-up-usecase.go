package usecases

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SignUpUsecaseInput struct {
	Name     string
	Email    string
	Password string
}

type SignUpUsecase struct {
	pgxPool     *pgxpool.Pool
	customerDAO daos.CustomerDAO
}

func NewSignUpUsecase(pgxPool *pgxpool.Pool, customerDAO daos.CustomerDAO) SignUpUsecase {
	return SignUpUsecase{pgxPool, customerDAO}
}

func (r SignUpUsecase) Execute(input SignUpUsecaseInput) error {
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

	tx, err := r.pgxPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	customerId := uuid.New()

	_, err = tx.Exec(context.Background(), "INSERT INTO customers (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)",
		customerId, input.Name, input.Email, string(hashedPassword), time.Now().UTC())
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO carts (id, customer_id, created_at) VALUES ($1, $2, $3)",
		uuid.New(), customerId, time.Now().UTC())
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
