package storage

import (
	"context"

	app "github.com/gerbenjacobs/go-webshop-course"
)

type BasketRepo struct {
	baskets map[int]app.Basket
}

func NewBasketRepo() *BasketRepo {
	return &BasketRepo{
		baskets: make(map[int]app.Basket),
	}
}

func (r *BasketRepo) GetBasket(ctx context.Context, userID int) (app.Basket, error) {
	basket, ok := r.baskets[userID]
	if !ok {
		basket = app.Basket{UserID: userID, Items: []app.BasketItem{}}
		r.baskets[userID] = basket
		return basket, nil
	}
	return basket, nil
}

func (r *BasketRepo) AddToBasket(ctx context.Context, userID, productID, quantity int) error {
	basket, ok := r.baskets[userID]
	if !ok {
		return app.ErrBasketNotFound
	}
	basket.Items = append(basket.Items, app.BasketItem{
		ProductID: productID,
		Quantity:  quantity,
	})
	r.baskets[userID] = basket
	return nil
}

func (r *BasketRepo) RemoveFromBasket(ctx context.Context, userID, productID, quantity int) error {
	basket, ok := r.baskets[userID]
	if !ok {
		return app.ErrBasketNotFound
	}
	for i, item := range basket.Items {
		if item.ProductID == productID {
			basket.Items = append(basket.Items[:i], basket.Items[i+1:]...)
			r.baskets[userID] = basket
			return nil
		}
	}
	// we have not found the product, but in essence that's the same as removing it
	// so we don't return an error here
	return nil
}
