package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

type OrderHandler struct {
	orderService     domain.OrderRepository
	orderItemService domain.OrderItemsRepository
}

func NewOrderHandler(orderService domain.OrderRepository, orderItemService domain.OrderItemsRepository) *OrderHandler {
	return &OrderHandler{
		orderService:     orderService,
		orderItemService: orderItemService,
	}
}

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	var input struct {
		UserId      string   `json:"user_id"`
		TotalAmount *float64 `json:"total_amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}

	if input.UserId == "" || input.TotalAmount == nil {
		response.Error(c, 400, "all fields are required")
		return
	}

	order := domain.Order{
		UserId:      input.UserId,
		TotalAmount: *input.TotalAmount,
	}

	err := h.orderService.CreateOrder(order)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}
	response.Success(c, 201, gin.H{"message": "order created"})
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
		response.Error(c, 500, "Something went wrong")
		return
	}
	response.Success(c, 200, gin.H{"order": order})
}

func (h *OrderHandler) UpdateStatusHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}
	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}

	if input.Status == "" {
		response.Error(c, 400, "status field is required")
		return
	}

	err := h.orderService.UpdateStatus(id, input.Status)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "order updated"})
}

func (h *OrderHandler) CreatItemsHandler(c *gin.Context) {
	var input struct {
		OrderID   string  `json:"order_id"`
		ProductID string  `json:"product_id"`
		Quantity  int     `json:"quantity"`
		UnitPrice float64 `json:"unit_price"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}

	if input.OrderID == "" || input.ProductID == "" || input.Quantity == 0 || input.UnitPrice == 0 {
		response.Error(c, 400, "all fields are required")
		return
	}
	orderItem := domain.OrderItem{
		OrderID:   input.OrderID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
		UnitPrice: input.UnitPrice,
	}

	err := h.orderItemService.CreatItems(orderItem)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}
	response.Success(c, 201, gin.H{"message": "order item created"})
}
func (h *OrderHandler) GetItemsbyOrderHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}

	orders, err := h.orderItemService.GetItemsbyOrder(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}
	response.Success(c, 200, gin.H{"orders": orders})
}
