package go_webshop_course

import "fmt"

type Product struct {
	ID          int
	Name        string
	Description string
	Image       string
	Price       float64
}

func (p Product) String() string {
	return fmt.Sprintf("[%d] %s - %s (€%.2f)", p.ID, p.Name, p.Description, p.Price)
}
