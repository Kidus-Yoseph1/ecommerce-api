package repository

import (
	"database/sql"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(order domain.Order) (*domain.Order, error) {
	err := r.db.QueryRow(`
        INSERT INTO orders (user_id, total_amount, status)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, total_amount, status, created_at, updated_at
    `, order.UserId, order.TotalAmount, order.Status).Scan(
		&order.Id,
		&order.UserId,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) GetOrder(id string) (*domain.Order, error) {
	order := &domain.Order{}

	err := r.db.QueryRow(`
	SELECT ID, user_id, total_amount, status, created_at, updated_at
	FROM orders
	WHERE id = $1
	`, id).Scan(
		&order.Id,
		&order.UserId,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepo) UpdateStatus(id string, status string) error {
	_, err := r.db.Exec(`
	UPDATE orders
	SET status = $1, updated_at = now()
	WHERE id = $2
	`, status, id,
	)
	return err
}

func (r *OrderRepo) CreatItems(item domain.OrderItem) error {
	_, err := r.db.Exec(`
	INSERT INTO order_items (order_id, product_id, quantity, unit_price)
	VALUES ($1, $2, $3, $4)
	`,
		item.OrderID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice)

	return err
}

func (r *OrderRepo) GetItemsbyOrder(id string) ([]domain.OrderItem, error) {
	order_items := []domain.OrderItem{}

	rows, err := r.db.Query(`
		SELECT id, order_id, product_id, quantity, unit_price, created_at
		FROM order_items
		WHERE order_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oi domain.OrderItem
		rows.Scan(
			&oi.ID,
			&oi.OrderID,
			&oi.ProductID,
			&oi.Quantity,
			&oi.UnitPrice,
			&oi.CreatedAt)
		order_items = append(order_items, oi)
	}
	return order_items, nil
}
