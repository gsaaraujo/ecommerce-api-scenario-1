package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddProductUsecaseInput struct {
	Name        string
	Description *string
	Price       int64
}

type AddProductUsecase struct {
	pgxPool *pgxpool.Pool
}

func NewAddProductUsecase(pgxPool *pgxpool.Pool) AddProductUsecase {
	return AddProductUsecase{pgxPool}
}

func (a *AddProductUsecase) Execute(input AddProductUsecaseInput) error {
	if input.Price == 0 {
		return errors.New("the product price cannot be zero")
	}

	tx, err := a.pgxPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	productId := uuid.New()
	_, err = tx.Exec(context.Background(), "INSERT INTO products (id, status, name, description, price, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		productId, "unpublished", input.Name, input.Description, input.Price, time.Now().UTC())
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO inventories (id, product_id, stock_quantity, created_at) VALUES ($1, $2, $3, $4)",
		uuid.New(), productId, 0, time.Now().UTC())
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
