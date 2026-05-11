package service

import (
	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type CartService struct {
	cartRepo     domain.CartRepository
	cartItemRepo domain.CartItemsRepository
}

func NewCartService(cartRepo domain.CartRepository, cartItemRepo domain.CartItemsRepository) *CartService {
	return &CartService{
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
	}
}

func (s *CartService) CreatCart(user_id string) error {
	cart := domain.Cart{
		UserId: user_id,
	}
	err := s.cartRepo.CreatCart(cart)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
	}
	return nil
}

func (s *CartService) GetCart(id string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCart(id)
	if err != nil {
		return nil, domain.ErrInternal("Something went wrong")
	}
	return cart, nil
}

func (s *CartService) ClearCart(id string) error {
	existing, err := s.cartRepo.GetCart(id)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
	}
	if existing == nil {
		return domain.ErrNotFound("Cart is empty")
	}

	err = s.cartRepo.ClearCart(id)
	if err != nil {
		return domain.ErrInternal("Something wnet wrong")
	}
	return nil
}

func (s *CartService) AddItem(cart_id string, product_id string, quantity int) error {
	item := domain.CartItem{
		CartID:    cart_id,
		ProductID: product_id,
		Quantity:  quantity,
	}

	err := s.cartItemRepo.AddItem(item)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
	}
	return nil
}

func (s *CartService) GetItems(id string) ([]domain.CartItem, error) {
	items, err := s.cartItemRepo.GetItems(id)
	if err != nil {
		return nil, domain.ErrInternal("Something went wrong")
	}
	if len(items) == 0 {
		return nil, nil
	}

	return items, nil
}

func (s *CartService) RemoveItem(id string) error {
	err := s.cartItemRepo.RemoveItem(id)
	if err != nil {
		return domain.ErrInternal("Somethig went wrong")
	}
	return nil
}
