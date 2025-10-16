package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DecreaseProductQuantityInCartUsecaseInput struct {
	CustomerId uuid.UUID
	ProductId  uuid.UUID
	Quantity   int32
}

type DecreaseProductQuantityInCartUsecase struct {
	pgxPool     *pgxpool.Pool
	cartDAO     daos.CartDAO
	cartItemDAO daos.CartItemDAO
}

func NewDecreaseProductQuantityInCartUsecase(pgxPool *pgxpool.Pool, cartDAO daos.CartDAO, cartItemDAO daos.CartItemDAO) DecreaseProductQuantityInCartUsecase {
	return DecreaseProductQuantityInCartUsecase{pgxPool, cartDAO, cartItemDAO}
}

func (i *DecreaseProductQuantityInCartUsecase) Execute(input DecreaseProductQuantityInCartUsecaseInput) error {
	if input.Quantity == 0 {
		return errors.New("you cannot decrease the quantity of product with a value equal to zero")
	}

	cartSchema, err := i.cartDAO.FindOneByCustomerId(input.CustomerId)
	if err != nil {
		return err
	}

	cartItemSchema, err := i.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)
	if err != nil {
		return err
	}

	if cartItemSchema == nil {
		return errors.New("product not found in cart")
	}

	if input.Quantity >= cartItemSchema.Quantity {
		_, err = i.pgxPool.Exec(context.Background(), "DELETE FROM cart_items WHERE id = $1", cartItemSchema.Id)
		return err
	}

	_, err = i.pgxPool.Exec(context.Background(), "UPDATE cart_items SET quantity = quantity - $1 WHERE id = $2",
		input.Quantity, cartItemSchema.Id)
	return err
}
