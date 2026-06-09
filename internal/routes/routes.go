package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kidus-yoseph1/ecommerce-api/internal/handler"
	"github.com/kidus-yoseph1/ecommerce-api/pkg/middleware"
)

func Setup(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	jwtSecret string,
) {
	// public
	r.POST("/auth/register", authHandler.RegisterHandler)
	r.POST("/auth/login", authHandler.LoginHandler)
	r.GET("/products", productHandler.ListProductsHandler)
	r.GET("/products/:id", productHandler.GetProductByIDHandler)
	r.POST("/webhooks/stripe", paymentHandler.WebhookHandler)

	// protected
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// cart
		protected.POST("/cart", cartHandler.CreateCartHandler)
		protected.GET("/cart", cartHandler.GetCartHandler)
		protected.GET("/cart/items", cartHandler.GetItemsHandler)
		protected.POST("/cart/items", cartHandler.AddItemHandler)
		protected.DELETE("/cart/items/:id", cartHandler.RemoveItemHandler)

		// orders
		protected.POST("/checkout", orderHandler.CreateOrderHandler)
		protected.GET("/orders/:id", orderHandler.GetOrderHandler)
		protected.GET("/orders/:id/items", orderHandler.GetItemsByOrderHandler)
	}

	// admin
	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/products", productHandler.AddProductHandler)
		admin.PUT("/products/:id", productHandler.UpdateProductHandler)
		admin.DELETE("/products/:id", productHandler.DeleteProductHandler)
	}
}
