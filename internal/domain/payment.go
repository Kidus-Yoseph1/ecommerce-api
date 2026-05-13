package domain

import "time"

type Payment struct {
	Id                    string    `json:"id"`
	OrderId               string    `json:"order_id"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
	Amount                float64   `json:"amount"`
	Status                string    `json:"status"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type PaymentRepository interface {
	CreatePayment(payment Payment) error
	UpdatePaymentStatus(stripePaymentIntentID string, status string) error
	GetPaymentByStripeID(stripeID string) (*Payment, error)
}

type PaymentServiceInterface interface {
	CreatePaymentIntent(orderID string, amount float64) (string, error)
	HandleWebhook(payload []byte, signature string, webhookSecret string) error
}
