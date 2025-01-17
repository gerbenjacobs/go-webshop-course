package storage

import (
	"context"
	"fmt"

	app "github.com/gerbenjacobs/go-webshop-course"
)

type ProductRepo struct {
	products map[int]app.Product
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		products: map[int]app.Product{
			1: {
				ID:          1,
				Name:        "Gopher plushie",
				Description: "A small purple Gophier plushie, perfect for kids and adults alike.",
				Image:       "",
				Price:       12.99,
			},
			2: {
				ID:          2,
				Name:        "PHP Elephant plushie",
				Description: "An elephant with the PHP logo, available in blue and pink",
				Image:       "",
				Price:       20,
			},
		},
	}
}

func (p *ProductRepo) GetAllProducts(_ context.Context) ([]app.Product, error) {
	var products []app.Product
	for _, product := range p.products {
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductRepo) GetProduct(ctx context.Context, productID int) (app.Product, error) {
	v, ok := p.products[productID]
	if !ok {
		return app.Product{}, fmt.Errorf("%w: for ID: %d", app.ErrProductNotFound, productID)
	}
	return v, nil
}
