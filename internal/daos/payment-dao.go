package daos

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (p *PaymentDAO) Create(paymentSchema PaymentSchema) error {
	_, err := p.pgxPool.Exec(context.Background(),
		"INSERT INTO payments (id, order_id, payment_gateway_name, payment_gateway_transaction_id, created_at) VALUES ($1, $2, $3, $4, $5)",
		paymentSchema.Id, paymentSchema.OrderId, paymentSchema.PaymentGatewayName, paymentSchema.PaymentGatewayTransactionId, paymentSchema.CreatedAt)

	return err
}

func (p *PaymentDAO) FindOneByCustomerId(customerId uuid.UUID) (*PaymentSchema, error) {
	var paymentSchema PaymentSchema

	err := p.pgxPool.QueryRow(context.Background(),
		"SELECT p.* FROM payments p JOIN orders o ON o.id = p.order_id WHERE o.customer_id = $1", customerId).
		Scan(&paymentSchema.Id, &paymentSchema.OrderId, &paymentSchema.PaymentGatewayName, &paymentSchema.PaymentGatewayTransactionId, &paymentSchema.CreatedAt)

	if err != nil && err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &paymentSchema, nil
}

func (p *PaymentDAO) DeletAll() error {
	_, err := p.pgxPool.Exec(context.Background(), "TRUNCATE TABLE payments CASCADE")
	return err
}
