package domain

import "time"

type Product struct {
	Id            string     `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Category      string     `json:"category"`
	Price         float64    `json:"price"`
	StockQuantity int        `json:"stock_quantity"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type ProductRepository interface {
	AddProduct(product Product) error
	GetProductbyId(id string) (*Product, error)
	ListProduct(category string) ([]Product, error)
	UpdateProduct(product Product) error
	DeleteProduct(id string) error
}

type ProductServiceInterface interface {
	AddProduct(name, description, category string, price float64, stockQuantity int) error
	GetProductByID(id string) (*Product, error)
	ListProducts(category string) ([]Product, error)
	UpdateProduct(id, name, description, category string, price float64, stockQuantity int) error
	DeleteProduct(id string) error
}
