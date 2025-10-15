package handlers

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type item struct {
	Id          uuid.UUID `json:"id"`
	ProductId   uuid.UUID `json:"productId"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Quantity    int32     `json:"quantity"`
	Price       int64     `json:"price"`
}

type GetCartHandlerOutput struct {
	CartId        uuid.UUID `json:"cartId"`
	TotalItems    int       `json:"totalItems"`
	TotalQuantity int32     `json:"totalQuantity"`
	TotalPrice    int64     `json:"totalPrice"`
	Items         []item    `json:"items"`
}

type GetCartHandler struct {
	pgxPool *pgxpool.Pool
	cartDAO daos.CartDAO
}

func NewGetCartHandler(pgxPool *pgxpool.Pool, cartDAO daos.CartDAO) GetCartHandler {
	return GetCartHandler{pgxPool, cartDAO}
}

func (g *GetCartHandler) Handle(c echo.Context) error {
	token := c.Get("customer").(*jwt.Token)
	claims := token.Claims.(*usecases.JwtAccessTokenClaims)

	cartSchema, err := g.cartDAO.FindOneByCustomerId(uuid.MustParse(claims.Subject))
	if err != nil {
		return err
	}

	if cartSchema == nil {
		return c.JSON(409, map[string]any{"message": "cart not found"})
	}

	rows, err := g.pgxPool.Query(context.Background(),
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
		`, claims.Subject)
	if err != nil {
		return err
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
			return err
		}

		records = append(records, item)
	}

	output := GetCartHandlerOutput{
		Items: []item{},
	}
	output.CartId = cartSchema.Id
	output.TotalItems = len(records)

	for _, record := range records {
		output.TotalQuantity += record.CartItemQuantity
		output.TotalPrice += record.ProductPrice * int64(record.CartItemQuantity)
		output.Items = append(output.Items, item{
			Id:          record.CartItemId,
			ProductId:   record.ProductId,
			Name:        record.ProductName,
			Description: record.ProductDescription,
			Quantity:    record.CartItemQuantity,
			Price:       record.ProductPrice,
		})
	}

	return c.JSON(200, map[string]any{"data": output})
}
