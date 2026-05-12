package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

type ProductHandler struct {
	productService domain.ProductRepository
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

	product := domain.Product{
		Name:          input.Name,
		Description:   input.Description,
		Category:      input.Category,
		Price:         input.Price,
		StockQuantity: input.StockQuantity,
	}

	err := h.productService.AddProduct(product)
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

func (h *ProductHandler) GetProductbyIdHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, 400, "id filed is required")
		return
	}
	product, err := h.productService.GetProductbyId(id)
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

func (h *ProductHandler) ListProductHandler(c *gin.Context) {
	category := c.Query("category")

	products, err := h.productService.ListProduct(category)
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

	product := domain.Product{
		Id:            id,
		Name:          input.Name,
		Description:   input.Description,
		Category:      input.Category,
		Price:         input.Price,
		StockQuantity: input.StockQuantity,
	}

	err := h.productService.UpdateProduct(product)
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
