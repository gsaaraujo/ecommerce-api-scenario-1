package usecases

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
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

	customerSchema := r.customerDAO.FindOneByEmail(input.Email)

	if customerSchema != nil {
		return errors.New("this email address has already been taken by someone")
	}

	hashedPassword := utils.GetOrThrow(bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost))

	tx := utils.GetOrThrow(r.pgxPool.Begin(context.Background()))

	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	customerId := uuid.New()

	_ = utils.GetOrThrow(tx.Exec(context.Background(), "INSERT INTO customers (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)",
		customerId, input.Name, input.Email, string(hashedPassword), time.Now().UTC()))

	_ = utils.GetOrThrow(tx.Exec(context.Background(), "INSERT INTO carts (id, customer_id, created_at) VALUES ($1, $2, $3)",
		uuid.New(), customerId, time.Now().UTC()))

	utils.ThrowOnError(tx.Commit(context.Background()))

	return nil
}
