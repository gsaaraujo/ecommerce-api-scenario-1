package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerSchema struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type CustomerDAO struct {
	pgxPool *pgxpool.Pool
}

func NewCustomerDAO(pgxPool *pgxpool.Pool) CustomerDAO {
	return CustomerDAO{pgxPool}
}

func (p *CustomerDAO) Create(customerSchema CustomerSchema) error {
	_, err := p.pgxPool.Exec(context.Background(),
		"INSERT INTO customers (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)",
		customerSchema.Id, customerSchema.Name, customerSchema.Email, customerSchema.Password, customerSchema.CreatedAt)

	return err
}

func (c *CustomerDAO) FindOneByEmail(email string) (*CustomerSchema, error) {
	var customerSchema CustomerSchema

	err := c.pgxPool.QueryRow(context.Background(),
		"SELECT id, name, email, password, created_at FROM customers WHERE email = $1", email).
		Scan(&customerSchema.Id, &customerSchema.Name, &customerSchema.Email, &customerSchema.Password, &customerSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &customerSchema, nil
}

func (c *CustomerDAO) DeletAll() error {
	_, err := c.pgxPool.Exec(context.Background(), "DELETE FROM customers")
	return err
}
