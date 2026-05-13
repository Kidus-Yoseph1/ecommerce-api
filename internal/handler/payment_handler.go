package handler

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

type PaymentHandler struct {
	paymentService domain.PaymentServiceInterface
	webhookSecret  string
}

func NewPaymentHandler(paymentService domain.PaymentServiceInterface, webhookSecret string) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		webhookSecret:  webhookSecret,
	}
}

func (h *PaymentHandler) WebhookHandler(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, 400, "could not read request body")
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		response.Error(c, 400, "missing stripe signature")
		return
	}

	err = h.paymentService.HandleWebhook(payload, signature, h.webhookSecret)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	// Stripe expects a 200 response — if you return anything else
	// it will retry the webhook
	response.Success(c, 200, gin.H{"received": true})
}
