package domain

import "time"

type Order struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderRepository interface {
	CreateOrder(order Order) error
	GetOrder(id string) (*Order, error)
	UpdateStatus(id string, status string) error
}

type OrderItemsRepository interface {
	CreatItems(item OrderItem) error
	GetItemsbyOrder(id string) ([]OrderItem, error)
}
