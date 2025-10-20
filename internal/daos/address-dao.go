package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressSchema struct {
	Id          uuid.UUID
	CustomerId  uuid.UUID
	IsDefault   bool
	Street      string
	City        string
	State       string
	Number      string
	ZipCode     string
	AddressLine string
	CreatedAt   time.Time
}

type AddressDAO struct {
	pgxPool *pgxpool.Pool
}

func NewAddressDAO(pgxPool *pgxpool.Pool) AddressDAO {
	return AddressDAO{pgxPool}
}

func (c *AddressDAO) Create(addressSchema AddressSchema) error {
	_, err := c.pgxPool.Exec(context.Background(),
		`INSERT INTO addresses (id, customer_id, is_default, street, city, state, number, zip_code, address_line, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		addressSchema.Id, addressSchema.CustomerId, addressSchema.IsDefault, addressSchema.Street, addressSchema.City, addressSchema.State,
		addressSchema.Number, addressSchema.ZipCode, addressSchema.AddressLine, addressSchema.CreatedAt)

	return err
}

func (c *AddressDAO) FindOneByIsDefault(isDefault bool) (bool, error) {
	var addressId uuid.UUID

	err := c.pgxPool.QueryRow(context.Background(), "SELECT id FROM addresses WHERE is_default = $1", isDefault).Scan(&addressId)

	if err != nil && err == pgx.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *AddressDAO) FindAllByCustomerId(customerId uuid.UUID) ([]AddressSchema, error) {
	rows, err := c.pgxPool.Query(context.Background(), "SELECT * FROM addresses WHERE customer_id = $1", customerId)
	if err != nil {
		return []AddressSchema{}, nil
	}

	var addressSchema []AddressSchema
	for rows.Next() {
		var item AddressSchema

		err := rows.Scan(&item.Id, &item.CustomerId, &item.IsDefault, &item.Street, &item.City,
			&item.State, &item.Number, &item.ZipCode, &item.AddressLine, &item.CreatedAt)
		if err != nil {
			return []AddressSchema{}, nil
		}

		addressSchema = append(addressSchema, item)
	}

	return addressSchema, nil
}

func (c *AddressDAO) DeletAll() error {
	_, err := c.pgxPool.Exec(context.Background(), "TRUNCATE TABLE addresses CASCADE")
	return err
}
