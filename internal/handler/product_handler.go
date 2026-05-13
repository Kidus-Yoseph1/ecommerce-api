package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

type ProductHandler struct {
	productService domain.ProductServiceInterface
}

func NewProductHandler(productService domain.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) AddProductHandler(c *gin.Context) {
	var input struct {
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		Category      string  `json:"category"`
		Price         float64 `json:"price"`
		StockQuantity int     `json:"stock_quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}

	if input.Name == "" || input.Category == "" || input.Price == 0 || input.StockQuantity == 0 {
		response.Error(c, 400, "all fields are required")
		return
	}

	err := h.productService.AddProduct(input.Name, input.Description, input.Category, input.Price, input.StockQuantity)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}
	response.Success(c, 201, gin.H{"message": "Product added"})
}

func (h *ProductHandler) GetProductByIDHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id filed is required")
		return
	}
	product, err := h.productService.GetProductByID(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}
	response.Success(c, 200, gin.H{"product": product})
}

func (h *ProductHandler) ListProductsHandler(c *gin.Context) {
	category := c.Query("category")

	products, err := h.productService.ListProducts(category)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}
	response.Success(c, 200, gin.H{"products": products})
}

func (h *ProductHandler) UpdateProductHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "if field is required")
		return
	}

	var input struct {
		Name          string  `json:"name"`
		Description   string  `json:"description"`
		Category      string  `json:"category"`
		Price         float64 `json:"price"`
		StockQuantity int     `json:"stock_quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}

	if input.Name == "" || input.Category == "" || input.Price == 0 || input.StockQuantity == 0 {
		response.Error(c, 400, "all fields are required")
		return
	}

	err := h.productService.UpdateProduct(id, input.Name, input.Description, input.Category, input.Price, input.StockQuantity)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "Product updated"})
}

func (h *ProductHandler) DeleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id field is required")
		return
	}
	err := h.productService.DeleteProduct(id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"message": "Product deleted"})
}
