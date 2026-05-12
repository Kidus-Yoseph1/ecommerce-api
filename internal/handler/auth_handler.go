package handler

import (
	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type AuthHandler struct {
	authService domain.AuthServiceInterface
}
