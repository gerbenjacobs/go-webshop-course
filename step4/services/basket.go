package services

import (
	"context"

	app "github.com/gerbenjacobs/go-webshop-course"
	"github.com/gerbenjacobs/go-webshop-course/storage"
)

type BasketSvc struct {
	repo storage.BasketRepository
}

func NewBasketService(repo storage.BasketRepository) *BasketSvc {
	return &BasketSvc{repo: repo}
}

func (b *BasketSvc) GetBasket(ctx context.Context, userID int) (app.Basket, error) {
	return b.repo.GetBasket(ctx, userID)
}
func (b *BasketSvc) AddToBasket(ctx context.Context, userID, productID, quantity int) error {
	return b.repo.AddToBasket(ctx, userID, productID, quantity)
}
func (b *BasketSvc) RemoveFromBasket(ctx context.Context, userID, productID, quantity int) error {
	return b.repo.RemoveFromBasket(ctx, userID, productID, quantity)
}
