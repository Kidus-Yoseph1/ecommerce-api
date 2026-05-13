package repository

import (
	"database/sql"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type CartRepo struct {
	db *sql.DB
}

func NewCartRepo(db *sql.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) CreatCart(cart domain.Cart) error {
	_, err := r.db.Exec(`
		INSERT INTO carts (user_id)
		VALUES ($1)
	`, cart.UserId)
	return err
}

func (r *CartRepo) GetCart(id string) (*domain.Cart, error) {
	cart := &domain.Cart{}

	err := r.db.QueryRow(`
		SELECT id, user_id, created_at, updated_at
		FROM carts
		WHERE id = $1
	`, id).Scan(
		&cart.Id,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepo) GetCartByUserID(userID string) (*domain.Cart, error) {
	cart := &domain.Cart{}

	err := r.db.QueryRow(`
		SELECT id, user_id, created_at, updated_at
		FROM carts
		WHERE user_id = $1
	`, userID).Scan(
		&cart.Id,
		&cart.UserId,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepo) ClearCart(cart_id string) error {
	_, err := r.db.Exec(`
		DELETE FROM cart_items WHERE cart_id = $1
	`, cart_id)
	return err
}

func (r *CartRepo) AddItem(item domain.CartItem) error {
	_, err := r.db.Exec(`
    INSERT INTO cart_items (cart_id, product_id, quantity)
    VALUES ($1, $2, $3)
    ON CONFLICT (cart_id, product_id)
    DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity,
                  updated_at = now()
	`, item.CartID, item.ProductID, item.Quantity,
	)

	return err
}

func (r *CartRepo) GetItems(cart_id string) ([]domain.CartItem, error) {
	items := []domain.CartItem{}

	rows, err := r.db.Query(`
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, p.price, ci.created_at, ci.updated_at
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		WHERE ci.cart_id = $1
	`, cart_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i domain.CartItem
		rows.Scan(&i.ID,
			&i.CartID,
			&i.ProductID,
			&i.Quantity,
			&i.Price,
			&i.CreatedAt,
			&i.UpdatedAt)

		items = append(items, i)
	}
	return items, nil
}

func (r *CartRepo) RemoveItem(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM cart_items WHERE id = $1
	`, id)
	return err
}
