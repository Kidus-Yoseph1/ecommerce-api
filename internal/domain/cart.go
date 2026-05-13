package domain

import "time"

type Cart struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItem struct {
	ID        string    `json:"id"`
	CartID    string    `json:"cart_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartRepository interface {
	CreatCart(cart Cart) error
	GetCart(id string) (*Cart, error)
	GetCartByUserID(userID string) (*Cart, error)
	ClearCart(id string) error
}

type CartItemsRepository interface {
	AddItem(item CartItem) error
	GetItems(cartId string) ([]CartItem, error)
	RemoveItem(id string) error
}

type CartServiceInterface interface {
	CreatCart(userID string) error
	GetCart(userID string) (*Cart, error)
	ClearCart(userID string) error
	AddItem(userID, productID string, quantity int) error
	GetItems(userID string) ([]CartItem, error)
	RemoveItem(id string) error
}
