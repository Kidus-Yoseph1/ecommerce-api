package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

type CartHandler struct {
	cartService domain.CartServiceInterface
}

func NewCartHandler(cartService domain.CartServiceInterface) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) CreateCartHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	err := h.cartService.CreatCart(userID)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 201, gin.H{"message": "cart created"})
}

func (h *CartHandler) GetCartHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"cart": cart})
}

func (h *CartHandler) ClearCartHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	err := h.cartService.ClearCart(userID)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "cart cleared"})
}

func (h *CartHandler) AddItemHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	var input struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}
	if input.ProductID == "" || input.Quantity == 0 {
		response.Error(c, 400, "all fields are required")
		return
	}

	err := h.cartService.AddItem(userID, input.ProductID, input.Quantity)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "item added"})
}

func (h *CartHandler) GetItemsHandler(c *gin.Context) {
	userID := c.GetString("user_id")

	items, err := h.cartService.GetItems(userID)
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

func (h *CartHandler) RemoveItemHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}

	err := h.cartService.RemoveItem(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "item removed"})
}
