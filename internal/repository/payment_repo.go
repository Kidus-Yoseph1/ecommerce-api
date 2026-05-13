package repository

import (
	"database/sql"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) CreatePayment(payment domain.Payment) error {
	_, err := r.db.Exec(`
		INSERT INTO payments (order_id, stripe_payment_intent_id, amount, status)
		VALUES ($1,$2,$3,$4)
	`, payment.OrderId,
		payment.StripePaymentIntentID,
		payment.Amount,
		payment.Status,
	)
	return err
}

func (r *PaymentRepo) GetPaymentByStripeID(stripeID string) (*domain.Payment, error) {
	payment := &domain.Payment{}

	err := r.db.QueryRow(`
        SELECT id, order_id, stripe_payment_intent_id, amount, status, created_at, updated_at
        FROM payments
        WHERE stripe_payment_intent_id = $1
    `, stripeID).Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.StripePaymentIntentID,
		&payment.Amount,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepo) UpdatePaymentStatus(id string, status string) error {
	_, err := r.db.Exec(`
		UPDATE payments
		SET status = $1, updated_at = now()
		WHERE id = $2
	`, status, id,
	)
	return err
}
