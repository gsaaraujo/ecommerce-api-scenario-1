package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartItemSchema struct {
	Id        uuid.UUID
	CartId    uuid.UUID
	ProductId uuid.UUID
	Quantity  int32
	CreatedAt time.Time
}

type CartItemDAO struct {
	pgxPool *pgxpool.Pool
}

func NewCartItemDAO(pgxPool *pgxpool.Pool) CartItemDAO {
	return CartItemDAO{pgxPool}
}

func (c *CartItemDAO) Create(cartItemSchema CartItemSchema) {
	_ = utils.GetOrThrow(c.pgxPool.Exec(context.Background(), "INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4, $5)",
		cartItemSchema.Id, cartItemSchema.CartId, cartItemSchema.ProductId, cartItemSchema.Quantity, cartItemSchema.CreatedAt))
}

func (c *CartItemDAO) FindAllByCartId(cartId uuid.UUID) []CartItemSchema {
	rows := utils.GetOrThrow(c.pgxPool.Query(context.Background(),
		"SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id = $1", cartId))

	var cartItemsSchema []CartItemSchema
	for rows.Next() {
		var item CartItemSchema

		utils.ThrowOnError(rows.Scan(&item.Id, &item.CartId, &item.ProductId, &item.Quantity, &item.CreatedAt))
		cartItemsSchema = append(cartItemsSchema, item)
	}

	return cartItemsSchema
}

func (c *CartItemDAO) ExistsByProductId(productId uuid.UUID) bool {
	var id uuid.UUID

	err := c.pgxPool.QueryRow(context.Background(), "SELECT id FROM cart_items WHERE product_id = $1", productId).Scan(&id)

	if err != nil && err == pgx.ErrNoRows {
		return false
	}

	if err != nil {
		panic(err)
	}

	return true
}

func (c *CartItemDAO) FindOneByCartIdAndProductId(cartId uuid.UUID, productId uuid.UUID) *CartItemSchema {
	var cartItemSchema CartItemSchema

	err := c.pgxPool.QueryRow(context.Background(),
		"SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id = $1 AND product_id = $2", cartId, productId).
		Scan(&cartItemSchema.Id, &cartItemSchema.CartId, &cartItemSchema.ProductId, &cartItemSchema.Quantity, &cartItemSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &cartItemSchema
}

func (c *CartItemDAO) DeletAll() {
	_ = utils.GetOrThrow(c.pgxPool.Exec(context.Background(), "TRUNCATE TABLE cart_items CASCADE"))
}
