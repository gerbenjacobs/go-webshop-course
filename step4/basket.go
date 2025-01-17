package go_webshop_course

import "errors"

var ErrBasketNotFound = errors.New("basket not found")

type Basket struct {
	UserID int
	Items  []BasketItem
}

type BasketItem struct {
	ProductID int
	Quantity  int
}
