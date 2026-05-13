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

func (s *CartService) CreatCart(userID string) error {
	cart := domain.Cart{
		UserId: userID,
	}
	err := s.cartRepo.CreatCart(cart)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}
	return nil
}

func (s *CartService) GetCart(userID string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, domain.ErrInternal("something went wrong")
	}
	if cart == nil {
		return nil, domain.ErrNotFound("cart not found")
	}
	return cart, nil
}

func (s *CartService) ClearCart(userID string) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}
	if cart == nil {
		return domain.ErrNotFound("cart not found")
	}

	err = s.cartRepo.ClearCart(cart.Id)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}
	return nil
}

func (s *CartService) AddItem(userID, productID string, quantity int) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}
	if cart == nil {
		return domain.ErrNotFound("cart not found")
	}

	item := domain.CartItem{
		CartID:    cart.Id,
		ProductID: productID,
		Quantity:  quantity,
	}

	err = s.cartItemRepo.AddItem(item)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}
	return nil
}

func (s *CartService) GetItems(userID string) ([]domain.CartItem, error) {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, domain.ErrInternal("something went wrong")
	}
	if cart == nil {
		return nil, domain.ErrNotFound("cart not found")
	}

	items, err := s.cartItemRepo.GetItems(cart.Id)
	if err != nil {
		return nil, domain.ErrInternal("something went wrong")
	}
	if len(items) == 0 {
		return nil, domain.ErrNotFound("cart is empty")
	}

	return items, nil
}

func (s *CartService) RemoveItem(id string) error {
	err := s.cartItemRepo.RemoveItem(id)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}
	return nil
}
