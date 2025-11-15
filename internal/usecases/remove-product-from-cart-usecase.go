package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
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
	cartSchema := r.cartDAO.FindOneByCustomerId(input.CustomerId)
	cartItemSchema := r.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)

	if cartItemSchema == nil {
		return errors.New("product not found in cart")
	}

	_ = utils.GetOrThrow(r.pgxPool.Exec(context.Background(), "DELETE FROM cart_items WHERE product_id = $1", input.ProductId))

	return nil
}
