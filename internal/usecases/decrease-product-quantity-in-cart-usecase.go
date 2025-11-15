package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
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

	cartSchema := i.cartDAO.FindOneByCustomerId(input.CustomerId)
	cartItemSchema := i.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)

	if cartItemSchema == nil {
		return errors.New("product not found in cart")
	}

	if input.Quantity >= cartItemSchema.Quantity {
		_ = utils.GetOrThrow(i.pgxPool.Exec(context.Background(), "DELETE FROM cart_items WHERE id = $1", cartItemSchema.Id))
	}

	_ = utils.GetOrThrow(i.pgxPool.Exec(context.Background(), "UPDATE cart_items SET quantity = quantity - $1 WHERE id = $2",
		input.Quantity, cartItemSchema.Id))

	return nil
}
