package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
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
	productExists := p.productDAO.ExistsById(input.ProductId)

	if !productExists {
		return errors.New("product not found")
	}

	_ = utils.GetOrThrow(p.pgxPool.Exec(context.Background(), "UPDATE products SET status = $1 WHERE id = $2", "published", input.ProductId))

	return nil
}
