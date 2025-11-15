package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentSchema struct {
	Id                          uuid.UUID
	OrderId                     uuid.UUID
	PaymentGatewayTransactionId string
	PaymentGatewayName          string
	CreatedAt                   time.Time
}

type PaymentDAO struct {
	pgxPool *pgxpool.Pool
}

func NewPaymentDAO(pgxPool *pgxpool.Pool) PaymentDAO {
	return PaymentDAO{pgxPool}
}

func (p *PaymentDAO) Create(paymentSchema PaymentSchema) {
	_ = utils.GetOrThrow(p.pgxPool.Exec(context.Background(),
		"INSERT INTO payments (id, order_id, payment_gateway_name, payment_gateway_transaction_id, created_at) VALUES ($1, $2, $3, $4, $5)",
		paymentSchema.Id, paymentSchema.OrderId, paymentSchema.PaymentGatewayName, paymentSchema.PaymentGatewayTransactionId, paymentSchema.CreatedAt))
}

func (p *PaymentDAO) FindOneByCustomerId(customerId uuid.UUID) *PaymentSchema {
	var paymentSchema PaymentSchema

	err := p.pgxPool.QueryRow(context.Background(),
		"SELECT p.* FROM payments p JOIN orders o ON o.id = p.order_id WHERE o.customer_id = $1", customerId).
		Scan(&paymentSchema.Id, &paymentSchema.OrderId, &paymentSchema.PaymentGatewayName, &paymentSchema.PaymentGatewayTransactionId, &paymentSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &paymentSchema
}

func (p *PaymentDAO) DeletAll() {
	_ = utils.GetOrThrow(p.pgxPool.Exec(context.Background(), "TRUNCATE TABLE payments CASCADE"))
}
