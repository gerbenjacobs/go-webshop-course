package services

import (
	"context"

	app "github.com/gerbenjacobs/go-webshop-course"
)

type ProductService interface {
	ListProducts(context.Context) ([]app.Product, error)
	ShowProduct(context.Context, int) (app.Product, error)
}

type BasketService interface {
	GetBasket(ctx context.Context, userID int) (app.Basket, error)
	AddToBasket(ctx context.Context, userID, productID, quantity int) error
	RemoveFromBasket(ctx context.Context, userID, productID, quantity int) error
}
