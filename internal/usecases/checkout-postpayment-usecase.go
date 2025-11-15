package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mercadopago/sdk-go/pkg/config"
)

type CheckoutPostpaymentUsecaseInput struct {
	CustomerId                  uuid.UUID
	AddressId                   uuid.UUID
	PaymentGatewayTransactionId string
}

type CheckoutPostpaymentUsecase struct {
	mercadoPagoConfig *config.Config
	pgxPool           *pgxpool.Pool
	cartDAO           daos.CartDAO
	cartItemDAO       daos.CartItemDAO
	inventoryDAO      daos.InventoryDAO
}

func NewCheckoutPostpaymentUsecase(mercadoPagoConfig *config.Config, pgxPool *pgxpool.Pool, cartDAO daos.CartDAO,
	cartItemDAO daos.CartItemDAO, inventoryDAO daos.InventoryDAO) CheckoutPostpaymentUsecase {
	return CheckoutPostpaymentUsecase{mercadoPagoConfig, pgxPool, cartDAO, cartItemDAO, inventoryDAO}
}

func (c *CheckoutPostpaymentUsecase) Execute(input CheckoutPostpaymentUsecaseInput) error {
	rows := utils.GetOrThrow(c.pgxPool.Query(context.Background(),
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
		`, input.CustomerId))

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

		utils.ThrowOnError(rows.Scan(&item.CartId, &item.CartItemId, &item.CartItemQuantity,
			&item.ProductId, &item.ProductName, &item.ProductDescription, &item.ProductPrice))

		records = append(records, item)
	}

	tx := utils.GetOrThrow(c.pgxPool.Begin(context.Background()))

	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	totalQuantity := int32(0)
	totalPrice := int64(0)

	for _, record := range records {
		totalQuantity += record.CartItemQuantity
		totalPrice += record.ProductPrice * int64(record.CartItemQuantity)

		_ = utils.GetOrThrow(tx.Exec(context.Background(),
			"UPDATE inventories SET stock_quantity = stock_quantity - $1 WHERE product_id = $2", record.CartItemQuantity, record.ProductId))
	}

	orderId := uuid.New()

	_ = utils.GetOrThrow(tx.Exec(context.Background(),
		"INSERT INTO orders (id, customer_id, total_price, total_quantity, created_at) VALUES ($1, $2, $3, $4, $5)",
		orderId, input.CustomerId, totalPrice, totalQuantity, time.Now().UTC()))

	for _, record := range records {
		_ = utils.GetOrThrow(tx.Exec(context.Background(),
			"INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
			uuid.New(), orderId, record.ProductId, record.CartItemQuantity, record.ProductPrice, time.Now().UTC()))
	}

	_ = utils.GetOrThrow(tx.Exec(context.Background(),
		"INSERT INTO payments (id, order_id, payment_gateway_name, payment_gateway_transaction_id, created_at) VALUES ($1, $2, $3, $4, $5)",
		uuid.New(), orderId, "mercado_pago", input.PaymentGatewayTransactionId, time.Now().UTC()))

	_ = utils.GetOrThrow(tx.Exec(context.Background(), "DELETE FROM cart_items WHERE cart_id = $1", records[0].CartId))
	utils.ThrowOnError(tx.Commit(context.Background()))

	return nil
}
