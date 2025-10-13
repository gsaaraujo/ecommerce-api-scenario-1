package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddProductToCartUsecaseInput struct {
	CustomerId uuid.UUID
	ProductId  uuid.UUID
	Quantity   int32
}

type AddProductToCartUsecase struct {
	pgxPool     *pgxpool.Pool
	cartDAO     daos.CartDAO
	cartItemDAO daos.CartItemDAO
	productDAO  daos.ProductDAO
}

func NewAddProductToCartUsecase(pgxPool *pgxpool.Pool, cartDAO daos.CartDAO, cartItemDAO daos.CartItemDAO, productDAO daos.ProductDAO) AddProductToCartUsecase {
	return AddProductToCartUsecase{pgxPool, cartDAO, cartItemDAO, productDAO}
}

func (a *AddProductToCartUsecase) Execute(input AddProductToCartUsecaseInput) error {
	if input.Quantity == 0 {
		return errors.New("product quantity cannot be zero")
	}

	productSchema, err := a.productDAO.FindOneById(input.ProductId)
	if err != nil {
		return err
	}

	if productSchema == nil {
		return errors.New("product not found")
	}

	cartSchema, err := a.cartDAO.FindOneByCustomerId(input.CustomerId)
	if err != nil {
		return err
	}

	if cartSchema != nil {
		cartItemSchema, err := a.cartItemDAO.FindOneByCartIdAndProductId(cartSchema.Id, input.ProductId)
		if err != nil {
			return err
		}

		if cartItemSchema != nil {
			_, err = a.pgxPool.Exec(context.Background(), "UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2",
				input.Quantity, cartItemSchema.Id)
			return err
		}

		_, err = a.pgxPool.Exec(context.Background(), "INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4, $5)",
			uuid.New(), cartSchema.Id, input.ProductId, input.Quantity, time.Now().UTC())
		return err
	}

	cartId := uuid.New()
	_, err = a.pgxPool.Exec(context.Background(), "INSERT INTO carts (id, customer_id, created_at) VALUES ($1, $2, $3)",
		cartId, input.CustomerId, time.Now().UTC())
	if err != nil {
		return err
	}

	_, err = a.pgxPool.Exec(context.Background(), "INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at) VALUES ($1, $2, $3, $4, $5)",
		uuid.New(), cartId, input.ProductId, input.Quantity, time.Now().UTC())
	return err

	// var cart *entities.Cart

	// cart, err = a.cartRepository.FindOneByCustomerId(input.CustomerId)
	// if err != nil {
	// 	return err
	// }

	// if cart == nil {
	// 	cart = utils.NewPointer(entities.NewCart(input.CustomerId))

	// 	err = cart.AddItem(productSchema.Id, input.Quantity, productSchema.Price)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = a.cartRepository.Create(*cart)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// err = cart.AddItem(productSchema.Id, input.Quantity, productSchema.Price)
	// if err != nil {
	// 	return err
	// }

	// return a.cartRepository.Update(*cart)
}
