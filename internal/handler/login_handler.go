package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/response"
)

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, 400, "invalid request body")
		return
	}

	if input.Email == "" || input.Password == "" {
		response.Error(c, 400, "all fields are required")
		return
	}

	token, err := h.authService.Login(input.Email, input.Password)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			response.Error(c, appErr.Code, appErr.Message)
			return
		}
		response.Error(c, 500, "Something went wrong")
		return
	}

	response.Success(c, 200, gin.H{"token": token})
}
