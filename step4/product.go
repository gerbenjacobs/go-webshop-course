package go_webshop_course

import (
	"errors"
	"fmt"
)

var ErrProductNotFound = errors.New("product not found")

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"desc"`
	Image       string  `json:"img"`
	Price       float64 `json:"price"`
}

func (p Product) String() string {
	return fmt.Sprintf("[%d] %s - %s (€%.2f)", p.ID, p.Name, p.Description, p.Price)
}

func (p Product) FormattedPrice() string {
	return fmt.Sprintf("€%.2f", p.Price)
}
