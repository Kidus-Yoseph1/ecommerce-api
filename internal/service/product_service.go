package service

import (
	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type ProductService struct {
	productRepo domain.ProductRepository
}

func NewProductService(ProductRepo domain.ProductRepository) *ProductService {
	return &ProductService{productRepo: ProductRepo}
}

func (s *ProductService) AddProduct(name string, description string, category string, price float64, stockquantity int) error {
	product := domain.Product{
		Name:          name,
		Description:   description,
		Category:      category,
		Price:         price,
		StockQuantity: stockquantity,
	}
	err := s.productRepo.AddProduct(product)
	if err != nil {
		return domain.ErrInternal("Product not created")
	}
	return nil
}

func (s *ProductService) GetProductbyId(id string) (*domain.Product, error) {
	product, err := s.productRepo.GetProductbyId(id)
	if err != nil {
		return nil, domain.ErrInternal("Something went wrong")
	}
	if product == nil {
		return nil, domain.ErrNotFound("Product not found")
	}
	return product, nil
}

func (s *ProductService) ListProduct(category string) ([]domain.Product, error) {
	products, err := s.productRepo.ListProduct(category)
	if err != nil {
		return nil, domain.ErrInternal("Something went wrong")
	}
	if len(products) == 0 {
		return nil, domain.ErrNotFound("Category not found")
	}

	return products, nil
}

func (s *ProductService) UpdateProduct(id string, name string, description string, category string, price float64, stockquantity int) error {
	productUpdated := domain.Product{
		Id:            id,
		Name:          name,
		Description:   description,
		Category:      category,
		Price:         price,
		StockQuantity: stockquantity,
	}
	err := s.productRepo.UpdateProduct(productUpdated)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
	}
	return nil
}

func (s *ProductService) DeleteProduct(id string) error {
	existing, err := s.productRepo.GetProductbyId(id)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
	}
	if existing == nil {
		return domain.ErrNotFound("Product does not exist")
	}

	err = s.productRepo.DeleteProduct(id)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
	}
	return nil
}
