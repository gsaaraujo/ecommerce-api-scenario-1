package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
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

func (p *InventoryDAO) Create(inventorySchema InventorySchema) {
	_ = utils.GetOrThrow(p.pgxPool.Exec(context.Background(),
		"INSERT INTO inventories (id, product_id, stock_quantity, created_at) VALUES ($1, $2, $3, $4)",
		inventorySchema.Id, inventorySchema.ProductId, inventorySchema.StockQuantity, inventorySchema.CreatedAt))
}

func (m *InventoryDAO) FindOneByProductId(productId uuid.UUID) *InventorySchema {
	var inventorySchema InventorySchema

	err := m.pgxPool.QueryRow(context.Background(),
		"SELECT id, product_id, stock_quantity, created_at FROM inventories WHERE product_id = $1", productId).
		Scan(&inventorySchema.Id, &inventorySchema.ProductId, &inventorySchema.StockQuantity, &inventorySchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &inventorySchema
}

func (m *InventoryDAO) ExistsById(id uuid.UUID) bool {
	var inventorySchema InventorySchema

	err := m.pgxPool.QueryRow(context.Background(), "SELECT id FROM inventories WHERE id = $1", id).
		Scan(&inventorySchema.Id)

	if err != nil && err == pgx.ErrNoRows {
		return false
	}

	if err != nil {
		panic(err)
	}

	return true
}

func (p *InventoryDAO) DeletAll() {
	_ = utils.GetOrThrow(p.pgxPool.Exec(context.Background(), "TRUNCATE TABLE inventories CASCADE"))
}
