package storage

import (
	"context"

	app "github.com/gerbenjacobs/go-webshop-course"
)

type ProductRepository interface {
	GetAllProducts(context.Context) ([]app.Product, error)
}
