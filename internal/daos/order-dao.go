package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderSchema struct {
	Id            uuid.UUID
	CustomerId    uuid.UUID
	TotalPrice    int64
	TotalQuantity int32
	CreatedAt     time.Time
}

type OrderDAO struct {
	pgxPool *pgxpool.Pool
}

func NewOrderDAO(pgxPool *pgxpool.Pool) OrderDAO {
	return OrderDAO{pgxPool}
}

func (o *OrderDAO) Create(orderSchema OrderSchema) error {
	_, err := o.pgxPool.Exec(context.Background(),
		"INSERT INTO carts (id, customer_id, total_price, total_quantity, created_at) VALUES ($1, $2, $3, $4, $5)",
		orderSchema.Id, orderSchema.CustomerId, orderSchema.TotalPrice, orderSchema.TotalQuantity, orderSchema.CreatedAt)

	return err
}

func (o *OrderDAO) FindOneByCustomerId(customerId uuid.UUID) (*OrderSchema, error) {
	var orderSchema OrderSchema

	err := o.pgxPool.QueryRow(context.Background(),
		"SELECT id, customer_id, total_price, total_quantity, created_at FROM orders WHERE customer_id = $1", customerId).
		Scan(&orderSchema.Id, &orderSchema.CustomerId, &orderSchema.TotalPrice, &orderSchema.TotalQuantity, &orderSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &orderSchema, nil
}

func (o *OrderDAO) DeletAll() error {
	_, err := o.pgxPool.Exec(context.Background(), "TRUNCATE TABLE orders CASCADE")
	return err
}
