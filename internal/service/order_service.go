package service

import (
	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type OrderService struct {
	orderRepo     domain.OrderRepository
	orderItemRepo domain.OrderItemsRepository
}

func NewOrderService(orderRepo domain.OrderRepository, orderItemRepo domain.OrderItemsRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (s *OrderService) CreateOrder(userId string, totalAmount float64) error {
	order := domain.Order{
		UserId:      userId,
		TotalAmount: totalAmount,
		Status:      "pending",
	}
	err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
	}
	return nil
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

func (s *OrderService) CreatItems(orderId string, productId string, quantity int, unitPrice float64) error {
	orderItem := domain.OrderItem{
		OrderID:   orderId,
		ProductID: productId,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
	err := s.orderItemRepo.CreatItems(orderItem)
	if err != nil {
		return domain.ErrInternal("Something went wrong")
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
