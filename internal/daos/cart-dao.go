package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
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

func (c *CartDAO) Create(cartSchema CartSchema) {
	_ = utils.GetOrThrow(c.pgxPool.Exec(context.Background(), "INSERT INTO carts (id, customer_id, created_at) VALUES ($1, $2, $3)",
		cartSchema.Id, cartSchema.CustomerId, cartSchema.CreatedAt))
}

func (c *CartDAO) FindOneByCustomerId(customerId uuid.UUID) *CartSchema {
	var cartSchema CartSchema

	err := c.pgxPool.QueryRow(context.Background(), "SELECT id, customer_id, created_at FROM carts WHERE customer_id = $1", customerId).
		Scan(&cartSchema.Id, &cartSchema.CustomerId, &cartSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &cartSchema
}

func (c *CartDAO) DeletAll() {
	_ = utils.GetOrThrow(c.pgxPool.Exec(context.Background(), "TRUNCATE TABLE carts CASCADE"))
}
