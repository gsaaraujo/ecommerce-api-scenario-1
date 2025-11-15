package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductSchema struct {
	Id          uuid.UUID
	Status      string
	Name        string
	Description *string
	Price       int64
	CreatedAt   time.Time
}

type ProductDAO struct {
	pgxPool *pgxpool.Pool
}

func NewProductDAO(pgxPool *pgxpool.Pool) ProductDAO {
	return ProductDAO{pgxPool}
}

func (p *ProductDAO) Create(productSchema ProductSchema) {
	_ = utils.GetOrThrow(p.pgxPool.Exec(context.Background(),
		"INSERT INTO products (id, status, name, description, price, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		productSchema.Id, productSchema.Status, productSchema.Name, productSchema.Description, productSchema.Price, productSchema.CreatedAt))
}

func (p *ProductDAO) FindOneById(id uuid.UUID) *ProductSchema {
	var productSchema ProductSchema

	err := p.pgxPool.QueryRow(context.Background(),
		"SELECT id, status, name, description, price, created_at FROM products WHERE id = $1", id).
		Scan(&productSchema.Id, &productSchema.Status, &productSchema.Name, &productSchema.Description, &productSchema.Price, &productSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &productSchema
}

func (p *ProductDAO) FindOneByName(name string) *ProductSchema {
	var productSchema ProductSchema

	err := p.pgxPool.QueryRow(context.Background(),
		"SELECT id, status, name, description, price, created_at FROM products WHERE name = $1", name).
		Scan(&productSchema.Id, &productSchema.Status, &productSchema.Name, &productSchema.Description, &productSchema.Price, &productSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &productSchema
}

func (p *ProductDAO) ExistsById(id uuid.UUID) bool {
	var productSchema ProductSchema

	err := p.pgxPool.QueryRow(context.Background(), "SELECT id FROM products WHERE id = $1", id).
		Scan(&productSchema.Id)

	if err != nil && err == pgx.ErrNoRows {
		return false
	}

	if err != nil {
		panic(err)
	}

	return true
}

func (p *ProductDAO) DeletAll() {
	_ = utils.GetOrThrow(p.pgxPool.Exec(context.Background(), "TRUNCATE TABLE products CASCADE"))
}
