package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncreaseProductQuantityInCartUsecaseInput struct {
	CustomerId uuid.UUID
	ProductId  uuid.UUID
	Quantity   int32
}

type IncreaseProductQuantityInCartUsecase struct {
	pgxPool     *pgxpool.Pool
	cartDAO     daos.CartDAO
	cartItemDAO daos.CartItemDAO
}

func NewIncreaseProductQuantityInCartUsecase(pgxPool *pgxpool.Pool, cartDAO daos.CartDAO, cartItemDAO daos.CartItemDAO) IncreaseProductQuantityInCartUsecase {
	return IncreaseProductQuantityInCartUsecase{pgxPool, cartDAO, cartItemDAO}
}

func (i *IncreaseProductQuantityInCartUsecase) Execute(input IncreaseProductQuantityInCartUsecaseInput) error {
	if input.Quantity == 0 {
		return errors.New("you cannot increase the quantity of product with a value equal to zero")
	}

	cartSchema, err := i.cartDAO.FindOneByCustomerId(input.CustomerId)
	if err != nil {
		return err
	}

	if cartSchema == nil {
		return errors.New("cart not found")
	}

	cartItemSchema, err := i.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)
	if err != nil {
		return err
	}

	if cartItemSchema == nil {
		return errors.New("product not found in cart")
	}

	_, err = i.pgxPool.Exec(context.Background(), "UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2",
		input.Quantity, cartItemSchema.Id)
	return err
}
