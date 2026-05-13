package service

import (
	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type OrderService struct {
	orderRepo     domain.OrderRepository
	orderItemRepo domain.OrderItemsRepository
	cartRepo      domain.CartRepository
	cartItemRepo  domain.CartItemsRepository
}

func NewOrderService(orderRepo domain.OrderRepository,
	orderItemRepo domain.OrderItemsRepository,
	cartRepo domain.CartRepository,
	cartItemRepo domain.CartItemsRepository,
) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		cartItemRepo:  cartItemRepo,
	}
}

func (s *OrderService) Checkout(userID string) (*domain.Order, error) {
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
		return nil, domain.ErrBadRequest("cart is empty")
	}

	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	order, err := s.orderRepo.CreateOrder(domain.Order{
		UserId:      userID,
		TotalAmount: total,
		Status:      "pending",
	})
	if err != nil {
		return nil, domain.ErrInternal("could not create order")
	}

	for _, item := range items {
		err = s.orderItemRepo.CreatItems(domain.OrderItem{
			OrderID:   order.Id,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.Price,
		})
		if err != nil {
			return nil, domain.ErrInternal("could not create order items")
		}
	}

	err = s.cartRepo.ClearCart(cart.Id)
	if err != nil {
		return nil, domain.ErrInternal("could not clear cart")
	}

	return order, nil
}

func (s *OrderService) GetOrder(id string) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrder(id)
	if err != nil {
		return nil, domain.ErrInternal("Something went wrong")
	}
	if order == nil {
		return nil, domain.ErrNotFound("Order not found")
	}

	return order, nil
}

func (s *OrderService) UpdateStatus(id string, status string) error {
	err := s.orderRepo.UpdateStatus(id, status)
	if err != nil {
		return domain.ErrInternal("something went wrong")
	}
	return nil
}

func (s *OrderService) GetItemsbyOrder(id string) ([]domain.OrderItem, error) {
	orderitems, err := s.orderItemRepo.GetItemsbyOrder(id)
	if err != nil {
		return nil, domain.ErrInternal("Something went wrong")
	}
	if len(orderitems) == 0 {
		return nil, domain.ErrNotFound("Order items not found")
	}
	return orderitems, nil
}
