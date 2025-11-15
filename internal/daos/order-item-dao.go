package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
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

func (o *OrderItemDAO) Create(orderItemSchema OrderItemSchema) {
	_ = utils.GetOrThrow(o.pgxPool.Exec(context.Background(),
		"INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		orderItemSchema.Id, orderItemSchema.OrderId, orderItemSchema.ProductId, orderItemSchema.Quantity, orderItemSchema.Price, orderItemSchema.CreatedAt))
}

func (o *OrderItemDAO) FindAllByOrderId(orderId uuid.UUID) []OrderItemSchema {
	rows := utils.GetOrThrow(o.pgxPool.Query(context.Background(),
		"SELECT id, order_id, product_id, quantity, price, created_at FROM order_items WHERE order_id = $1", orderId))

	var cartItemsSchema []OrderItemSchema
	for rows.Next() {
		var item OrderItemSchema

		utils.ThrowOnError(rows.Scan(&item.Id, &item.OrderId, &item.ProductId, &item.Quantity, &item.Price, &item.CreatedAt))
		cartItemsSchema = append(cartItemsSchema, item)
	}

	return cartItemsSchema
}

func (o *OrderItemDAO) DeletAll() {
	_ = utils.GetOrThrow(o.pgxPool.Exec(context.Background(), "TRUNCATE TABLE order_items CASCADE"))
}
