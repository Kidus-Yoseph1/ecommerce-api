package service

import (
	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

type PaymentService struct {
	paymentRepo domain.PaymentRepository
	orderRepo   domain.OrderRepository
	cartRepo    domain.CartRepository
}

func NewPaymentService(
	paymentRepo domain.PaymentRepository,
	orderRepo domain.OrderRepository,
	cartRepo domain.CartRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
	}
}

func (s *PaymentService) CreatePaymentIntent(orderID string, amount float64) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", domain.ErrInternal("could not create payment intent")
	}

	payment := domain.Payment{
		OrderId:               orderID,
		StripePaymentIntentID: pi.ID,
		Amount:                amount,
		Status:                "pending",
	}

	err = s.paymentRepo.CreatePayment(payment)
	if err != nil {
		return "", domain.ErrInternal("could not save payment")
	}

	return pi.ClientSecret, nil
}

func (s *PaymentService) HandleWebhook(payload []byte, signature string, webhookSecret string) error {
	event, err := webhook.ConstructEventWithOptions(payload, signature, webhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)

	if err != nil {
		return domain.ErrBadRequest("invalid webhook signature")
	}

	switch event.Type {
	case "payment_intent.succeeded":
		pi, ok := event.Data.Object["id"].(string)
		if !ok {
			return domain.ErrInternal("could not parse payment intent")
		}

		err = s.paymentRepo.UpdatePaymentStatus(pi, "succeeded")
		if err != nil {
			return domain.ErrInternal("could not update payment")
		}

		payment, err := s.paymentRepo.GetPaymentByStripeID(pi)
		if err != nil || payment == nil {
			return domain.ErrInternal("could not find payment")
		}

		err = s.orderRepo.UpdateStatus(payment.OrderId, "paid")
		if err != nil {
			return domain.ErrInternal("could not update order")
		}

		order, err := s.orderRepo.GetOrder(payment.OrderId)
		if err != nil || order == nil {
			return domain.ErrInternal("could not find order")
		}

		cart, err := s.cartRepo.GetCartByUserID(order.UserId)
		if err != nil || cart == nil {
			return domain.ErrInternal("could not find cart")
		}

		err = s.cartRepo.ClearCart(cart.Id)
		if err != nil {
			return domain.ErrInternal("could not clear cart")
		}

	case "payment_intent.payment_failed":
		pi, ok := event.Data.Object["id"].(string)
		if !ok {
			return domain.ErrInternal("could not parse payment intent")
		}

		err = s.paymentRepo.UpdatePaymentStatus(pi, "failed")
		if err != nil {
			return domain.ErrInternal("could not update payment")
		}
	}

	return nil
}
