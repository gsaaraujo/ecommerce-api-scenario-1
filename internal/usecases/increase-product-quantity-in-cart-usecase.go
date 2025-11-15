package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncreaseProductQuantityInCartUsecaseInput struct {
	CustomerId uuid.UUID
	ProductId  uuid.UUID
	Quantity   int32
}

type IncreaseProductQuantityInCartUsecase struct {
	pgxPool      *pgxpool.Pool
	cartDAO      daos.CartDAO
	cartItemDAO  daos.CartItemDAO
	inventoryDAO daos.InventoryDAO
}

func NewIncreaseProductQuantityInCartUsecase(pgxPool *pgxpool.Pool, cartDAO daos.CartDAO,
	cartItemDAO daos.CartItemDAO, inventoryDAO daos.InventoryDAO) IncreaseProductQuantityInCartUsecase {
	return IncreaseProductQuantityInCartUsecase{pgxPool, cartDAO, cartItemDAO, inventoryDAO}
}

func (i *IncreaseProductQuantityInCartUsecase) Execute(input IncreaseProductQuantityInCartUsecaseInput) error {
	if input.Quantity == 0 {
		return errors.New("you cannot increase the quantity of product with a value equal to zero")
	}

	cartSchema := i.cartDAO.FindOneByCustomerId(input.CustomerId)
	cartItemSchema := i.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)

	if cartItemSchema == nil {
		return errors.New("product not found in cart")
	}

	inventorySchema := i.inventoryDAO.FindOneByProductId(input.ProductId)

	if input.Quantity > inventorySchema.StockQuantity {
		return errors.New("product quantity exceeds the stock available")
	}

	_ = utils.GetOrThrow(i.pgxPool.Exec(context.Background(), "UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2",
		input.Quantity, cartItemSchema.Id))

	return nil
}
