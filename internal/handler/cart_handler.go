package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

type CartHandler struct {
	cartService     domain.CartRepository
	cartItemService domain.CartItemsRepository
}

func NewCartHandler(cartService domain.CartRepository, cartItemService domain.CartItemsRepository) *CartHandler {
	return &CartHandler{
		cartService:     cartService,
		cartItemService: cartItemService,
	}
}

func (h *CartHandler) CreateCartHandler(c *gin.Context) {
	var input struct {
		UserId string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request field")
		return
	}
	if input.UserId == "" {
		response.Error(c, 400, "all fields required")
		return
	}

	cart := domain.Cart{
		UserId: input.UserId,
	}

	err := h.cartService.CreatCart(cart)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 201, gin.H{"message": "product created"})
}

func (h *CartHandler) GetCartHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}

	cart, err := h.cartService.GetCart(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"cart": cart})

}

func (h *CartHandler) ClearCartHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}

	err := h.cartService.ClearCart(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "cart cleared"})
}

func (h *CartHandler) AddItemHandler(c *gin.Context) {
	var input struct {
		CartID    string `json:"cart_id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}
	if input.CartID == "" || input.ProductID == "" || input.Quantity == 0 {
		response.Error(c, 400, "all fields are required")
		return
	}

	cartItem := domain.CartItem{
		CartID:    input.CartID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	err := h.cartItemService.AddItem(cartItem)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "item added"})
}

func (h *CartHandler) GetItemsHandler(c *gin.Context) {
	cartId := c.Param("id")
	if cartId == "" {
		response.Error(c, 400, "id field is required")
		return
	}

	items, err := h.cartItemService.GetItems(cartId)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
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

	err := h.cartItemService.RemoveItem(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "item removed"})
}
