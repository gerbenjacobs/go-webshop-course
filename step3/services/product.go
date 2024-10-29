package services

import (
	"context"

	app "github.com/gerbenjacobs/go-webshop-course"
	"github.com/gerbenjacobs/go-webshop-course/storage"
)

type ProductSvc struct {
	repo storage.ProductRepository
}

func NewProductService(repo storage.ProductRepository) *ProductSvc {
	return &ProductSvc{repo: repo}
}

func (p *ProductSvc) ListProducts(ctx context.Context) ([]app.Product, error) {
	return p.repo.GetAllProducts(ctx)
}
