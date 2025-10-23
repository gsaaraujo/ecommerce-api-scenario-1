package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type CheckoutPrepaymentInput struct {
	CustomerId uuid.UUID
}

type CheckoutPrepaymentOutput struct {
	PreferenceId string
}

type CheckoutPrepayment struct {
	pgxPool          *pgxpool.Pool
	preferenceClient preference.Client
}

func NewCheckoutPrepayment(pgxPool *pgxpool.Pool, preferenceClient preference.Client) CheckoutPrepayment {
	return CheckoutPrepayment{pgxPool, preferenceClient}
}

func (c *CheckoutPrepayment) Execute(input CheckoutPrepaymentInput) (CheckoutPrepaymentOutput, error) {
	rows, err := c.pgxPool.Query(context.Background(),
		`
			SELECT
				c.id AS cart_id,
				ci.id AS cart_item_id,
				ci.quantity AS cart_item_quantity,
				p.id AS product_id,
				p.name AS product_name,
				p.description AS product_description,
				p.price AS product_price
			FROM carts c
			JOIN cart_items ci
				ON ci.cart_id = c.id
			JOIN products p
				ON ci.product_id = p.id
			WHERE c.customer_id = $1
		`, input.CustomerId)
	if err != nil {
		return CheckoutPrepaymentOutput{}, err
	}

	type schema struct {
		CartId             uuid.UUID
		CartItemId         uuid.UUID
		ProductId          uuid.UUID
		CartItemQuantity   int32
		ProductName        string
		ProductDescription *string
		ProductPrice       int64
	}

	records := []schema{}
	for rows.Next() {
		var item schema
		err := rows.Scan(&item.CartId, &item.CartItemId, &item.CartItemQuantity,
			&item.ProductId, &item.ProductName, &item.ProductDescription, &item.ProductPrice)
		if err != nil {
			return CheckoutPrepaymentOutput{}, err
		}

		records = append(records, item)
	}

	itemsRequest := []preference.ItemRequest{}
	for _, record := range records {
		itemsRequest = append(itemsRequest, preference.ItemRequest{
			ID:          record.ProductId.String(),
			Title:       record.ProductName,
			Description: *record.ProductDescription,
			Quantity:    int(record.CartItemQuantity),
			UnitPrice:   float64(record.ProductPrice),
		})
	}

	preferenceResponse, err := c.preferenceClient.Create(context.Background(), preference.Request{
		Items: itemsRequest,
	})
	if err != nil {
		return CheckoutPrepaymentOutput{}, err
	}

	return CheckoutPrepaymentOutput{
		PreferenceId: preferenceResponse.ID,
	}, nil
}
