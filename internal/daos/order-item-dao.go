package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItemSchema struct {
	Id        uuid.UUID
	OrderId   uuid.UUID
	ProductId uuid.UUID
	Quantity  int32
	Price     int64
	CreatedAt time.Time
}

type OrderItemDAO struct {
	pgxPool *pgxpool.Pool
}

func NewOrderItemDAO(pgxPool *pgxpool.Pool) OrderItemDAO {
	return OrderItemDAO{pgxPool}
}

func (o *OrderItemDAO) Create(orderItemSchema OrderItemSchema) error {
	_, err := o.pgxPool.Exec(context.Background(),
		"INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		orderItemSchema.Id, orderItemSchema.OrderId, orderItemSchema.ProductId, orderItemSchema.Quantity, orderItemSchema.Price, orderItemSchema.CreatedAt)

	return err
}

func (o *OrderItemDAO) FindAllByOrderId(orderId uuid.UUID) ([]OrderItemSchema, error) {
	rows, err := o.pgxPool.Query(context.Background(),
		"SELECT id, order_id, product_id, quantity, price, created_at FROM order_items WHERE order_id = $1", orderId)
	if err != nil {
		return []OrderItemSchema{}, nil
	}

	var cartItemsSchema []OrderItemSchema
	for rows.Next() {
		var item OrderItemSchema

		err := rows.Scan(&item.Id, &item.OrderId, &item.ProductId, &item.Quantity, &item.Price, &item.CreatedAt)
		if err != nil {
			return []OrderItemSchema{}, nil
		}

		cartItemsSchema = append(cartItemsSchema, item)
	}

	return cartItemsSchema, nil
}

func (o *OrderItemDAO) DeletAll() error {
	_, err := o.pgxPool.Exec(context.Background(), "TRUNCATE TABLE order_items CASCADE")
	return err
}
