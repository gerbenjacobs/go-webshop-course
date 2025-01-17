package storage

import (
	"context"

	app "github.com/gerbenjacobs/go-webshop-course"
)

type ProductRepository interface {
	GetAllProducts(context.Context) ([]app.Product, error)
	GetProduct(ctx context.Context, productID int) (app.Product, error)
}

type BasketRepository interface {
	GetBasket(ctx context.Context, userID int) (app.Basket, error)
	AddToBasket(ctx context.Context, userID, productID, quantity int) error
	RemoveFromBasket(ctx context.Context, userID, productID, quantity int) error
}
