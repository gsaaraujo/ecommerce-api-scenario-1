package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RemoveProductFromCartUsecaseInput struct {
	CustomerId uuid.UUID
	ProductId  uuid.UUID
}

type RemoveProductFromCartUsecase struct {
	pgxPool     *pgxpool.Pool
	cartDAO     daos.CartDAO
	cartItemDAO daos.CartItemDAO
}

func NewRemoveProductFromCartUsecase(pgxPool *pgxpool.Pool, cartDAO daos.CartDAO, cartItemDAO daos.CartItemDAO) RemoveProductFromCartUsecase {
	return RemoveProductFromCartUsecase{pgxPool, cartDAO, cartItemDAO}
}

func (r *RemoveProductFromCartUsecase) Execute(input RemoveProductFromCartUsecaseInput) error {
	cartSchema, err := r.cartDAO.FindOneByCustomerId(input.CustomerId)
	if err != nil {
		return err
	}

	cartItemSchema, err := r.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)
	if err != nil {
		return err
	}

	if cartItemSchema == nil {
		return errors.New("product not found in cart")
	}

	_, err = r.pgxPool.Exec(context.Background(), "DELETE FROM cart_items WHERE product_id = $1", input.ProductId)
	return err
}
