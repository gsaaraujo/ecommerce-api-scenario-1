package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PublishProductUsecaseInput struct {
	ProductId uuid.UUID
}

type PublishProductUsecase struct {
	pgxPool    *pgxpool.Pool
	productDAO daos.ProductDAO
}

func NewPublishProductUsecase(pgxPool *pgxpool.Pool, productDAO daos.ProductDAO) PublishProductUsecase {
	return PublishProductUsecase{pgxPool, productDAO}
}

func (p *PublishProductUsecase) Execute(input PublishProductUsecaseInput) error {
	productExists, err := p.productDAO.ExistsById(input.ProductId)
	if err != nil {
		return err
	}

	if !productExists {
		return errors.New("product not found")
	}

	_, err = p.pgxPool.Exec(context.Background(), "UPDATE products SET status = $1 WHERE id = $2", "published", input.ProductId)
	if err != nil {
		return err
	}

	return nil
}
