package domain

import "time"

type Payment struct {
	Id                    string    `json:"id"`
	OrderId               string    `json:"user_id"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
	Amount                float64   `json:"amount"`
	Status                string    `json:"status"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type PaymentRepository interface {
	CreatePayment(payment Payment) error
	UpdatePaymentStatus(id string, status string) error
}
