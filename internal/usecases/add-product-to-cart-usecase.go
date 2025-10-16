package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddProductToCartUsecaseInput struct {
	CustomerId uuid.UUID
	ProductId  uuid.UUID
	Quantity   int32
}

type AddProductToCartUsecase struct {
	pgxPool      *pgxpool.Pool
	cartDAO      daos.CartDAO
	cartItemDAO  daos.CartItemDAO
	productDAO   daos.ProductDAO
	inventoryDAO daos.InventoryDAO
}

func NewAddProductToCartUsecase(pgxPool *pgxpool.Pool, cartDAO daos.CartDAO, cartItemDAO daos.CartItemDAO,
	productDAO daos.ProductDAO, inventoryDAO daos.InventoryDAO) AddProductToCartUsecase {
	return AddProductToCartUsecase{pgxPool, cartDAO, cartItemDAO, productDAO, inventoryDAO}
}

func (a *AddProductToCartUsecase) Execute(input AddProductToCartUsecaseInput) error {
	if input.Quantity == 0 {
		return errors.New("product quantity cannot be zero")
	}

	productSchema, err := a.productDAO.FindOneById(input.ProductId)
	if err != nil {
		return err
	}

	if productSchema == nil {
		return errors.New("product not found")
	}

	inventorySchema, err := a.inventoryDAO.FindOneByProductId(input.ProductId)
	if err != nil {
		return err
	}

	if input.Quantity > inventorySchema.StockQuantity {
		return errors.New("product quantity exceeds the stock available")
	}

	cartSchema, err := a.cartDAO.FindOneByCustomerId(input.CustomerId)
	if err != nil {
		return err
	}

	cartItemSchema, err := a.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)
	if err != nil {
		return err
	}

	if cartItemSchema != nil {
		_, err = a.pgxPool.Exec(context.Background(), "UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2",
			input.Quantity, cartItemSchema.Id)
		return err
	}

	_, err = a.pgxPool.Exec(context.Background(), "INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4, $5)",
		uuid.New(), cartSchema.Id, input.ProductId, input.Quantity, time.Now().UTC())
	return err
}
