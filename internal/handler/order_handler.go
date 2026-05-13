package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

type OrderHandler struct {
	orderService domain.OrderServiceInterface
}

func NewOrderHandler(orderService domain.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	order, clientSecret, err := h.orderService.Checkout(userID)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 201, gin.H{
		"order_id":      order.Id,
		"client_secret": clientSecret,
	})
}

func (h *OrderHandler) GetOrderHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"order": order})
}

func (h *OrderHandler) GetItemsByOrderHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}

	items, err := h.orderService.GetItemsbyOrder(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"items": items})
}
