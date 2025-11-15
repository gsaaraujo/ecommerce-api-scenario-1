package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddStockUsecaseInput struct {
	InventoryId uuid.UUID
	Stock       int32
}

type AddStockUsecase struct {
	pgxPool      *pgxpool.Pool
	inventoryDAO daos.InventoryDAO
}

func NewAddStockUsecase(pgxPool *pgxpool.Pool, inventoryDAO daos.InventoryDAO) AddStockUsecase {
	return AddStockUsecase{pgxPool, inventoryDAO}
}

func (a *AddStockUsecase) Execute(input AddStockUsecaseInput) error {
	if input.Stock == 0 {
		return errors.New("stock quantity must be higher than zero")
	}

	inventoryExists := a.inventoryDAO.ExistsById(input.InventoryId)

	if !inventoryExists {
		return errors.New("inventory not found")
	}

	_ = utils.GetOrThrow(a.pgxPool.Exec(context.Background(),
		"UPDATE inventories SET stock_quantity = stock_quantity + $1 WHERE id = $2", input.Stock, input.InventoryId))

	return nil
}
