package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventorySchema struct {
	Id            uuid.UUID
	ProductId     uuid.UUID
	StockQuantity int32
	CreatedAt     time.Time
}

type InventoryDAO struct {
	pgxPool *pgxpool.Pool
}

func NewInventoryDAO(pgxPool *pgxpool.Pool) InventoryDAO {
	return InventoryDAO{pgxPool}
}

func (p *InventoryDAO) Create(inventorySchema InventorySchema) error {
	_, err := p.pgxPool.Exec(context.Background(),
		"INSERT INTO inventories (id, product_id, stock_quantity, created_at) VALUES ($1, $2, $3, $4)",
		inventorySchema.Id, inventorySchema.ProductId, inventorySchema.StockQuantity, inventorySchema.CreatedAt)

	return err
}

func (m *InventoryDAO) FindOneByProductId(productId uuid.UUID) (*InventorySchema, error) {
	var inventorySchema InventorySchema

	err := m.pgxPool.QueryRow(context.Background(),
		"SELECT id, product_id, stock_quantity, created_at FROM inventories WHERE product_id = $1", productId).
		Scan(&inventorySchema.Id, &inventorySchema.ProductId, &inventorySchema.StockQuantity, &inventorySchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &inventorySchema, nil
}

func (m *InventoryDAO) ExistsById(id uuid.UUID) (bool, error) {
	var inventorySchema InventorySchema

	err := m.pgxPool.QueryRow(context.Background(), "SELECT id FROM inventories WHERE id = $1", id).
		Scan(&inventorySchema.Id)

	if err != nil && err == pgx.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *InventoryDAO) DeletAll() error {
	_, err := p.pgxPool.Exec(context.Background(), "TRUNCATE TABLE inventories CASCADE")
	return err
}
