package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartSchema struct {
	Id         uuid.UUID
	CustomerId uuid.UUID
	CreatedAt  time.Time
}

type CartDAO struct {
	pgxPool *pgxpool.Pool
}

func NewCartDAO(pgxPool *pgxpool.Pool) CartDAO {
	return CartDAO{pgxPool}
}

func (c *CartDAO) Create(cartSchema CartSchema) error {
	_, err := c.pgxPool.Exec(context.Background(), "INSERT INTO carts (id, customer_id, created_at) VALUES ($1, $2, $3)",
		cartSchema.Id, cartSchema.CustomerId, cartSchema.CreatedAt)

	return err
}

func (c *CartDAO) FindOneByCustomerId(customerId uuid.UUID) (*CartSchema, error) {
	var cartSchema CartSchema

	err := c.pgxPool.QueryRow(context.Background(), "SELECT id, customer_id, created_at FROM carts WHERE customer_id = $1", customerId).
		Scan(&cartSchema.Id, &cartSchema.CustomerId, &cartSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &cartSchema, nil
}

func (c *CartDAO) DeletAll() error {
	_, err := c.pgxPool.Exec(context.Background(), "TRUNCATE TABLE carts CASCADE")
	return err
}
