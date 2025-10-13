package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (c *CartItemDAO) Create(cartItemSchema CartItemSchema) error {
	_, err := c.pgxPool.Exec(context.Background(), "INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4, $5)",
		cartItemSchema.Id, cartItemSchema.CartId, cartItemSchema.ProductId, cartItemSchema.Quantity, cartItemSchema.CreatedAt)

	return err
}

func (c *CartItemDAO) FindAllByCartId(cartId uuid.UUID) ([]CartItemSchema, error) {
	rows, err := c.pgxPool.Query(context.Background(),
		"SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id = $1", cartId)
	if err != nil {
		return []CartItemSchema{}, nil
	}

	var cartItemsSchema []CartItemSchema
	for rows.Next() {
		var item CartItemSchema

		err := rows.Scan(&item.Id, &item.CartId, &item.ProductId, &item.Quantity, &item.CreatedAt)
		if err != nil {
			return []CartItemSchema{}, nil
		}

		cartItemsSchema = append(cartItemsSchema, item)
	}

	return cartItemsSchema, nil
}

func (c *CartItemDAO) ExistsByProductId(productId uuid.UUID) (bool, error) {
	var id uuid.UUID

	err := c.pgxPool.QueryRow(context.Background(), "SELECT id FROM cart_items WHERE product_id = $1", productId).Scan(&id)
	if err != nil && err == pgx.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *CartItemDAO) FindOneByCartIdAndProductId(cartId uuid.UUID, productId uuid.UUID) (*CartItemSchema, error) {
	var cartItemSchema CartItemSchema

	err := c.pgxPool.QueryRow(context.Background(),
		"SELECT id, cart_id, product_id, quantity, created_at FROM cart_items WHERE cart_id = $1 AND product_id = $2", cartId, productId).
		Scan(&cartItemSchema.Id, &cartItemSchema.CartId, &cartItemSchema.ProductId, &cartItemSchema.Quantity, &cartItemSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &cartItemSchema, nil
}

func (c *CartItemDAO) DeletAll() error {
	_, err := c.pgxPool.Exec(context.Background(), "TRUNCATE TABLE cart_items CASCADE")
	return err
}
