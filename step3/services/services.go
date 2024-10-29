package services

import (
	"context"

	app "github.com/gerbenjacobs/go-webshop-course"
)

type ProductService interface {
	ListProducts(context.Context) ([]app.Product, error)
}
