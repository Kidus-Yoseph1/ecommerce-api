package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kidus-yoseph1/ecommerce-api/internal/db"
	"github.com/kidus-yoseph1/ecommerce-api/internal/handler"
	"github.com/kidus-yoseph1/ecommerce-api/internal/repository"
	"github.com/kidus-yoseph1/ecommerce-api/internal/routes"
	"github.com/kidus-yoseph1/ecommerce-api/internal/service"
	stripe "github.com/stripe/stripe-go/v76"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	database, err := db.Connect()
	if err != nil {
		log.Fatal("could not connect to database:", err)
	}
	defer database.Close()
	log.Println("database connected successfully")

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// repositories
	userRepo := repository.NewUserRepo(database)
	productRepo := repository.NewProductRepo(database)
	cartRepo := repository.NewCartRepo(database)
	orderRepo := repository.NewOrderRepo(database)
	paymentRepo := repository.NewPaymentRepo(database)

	// services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, cartRepo)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo, cartRepo)
	orderService := service.NewOrderService(orderRepo, orderRepo, cartRepo, cartRepo, paymentService)

	// handlers
	authHandler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService, os.Getenv("STRIPE_WEBHOOK_SECRET"))

	// router
	r := gin.Default()
	routes.Setup(r, authHandler, productHandler, cartHandler, orderHandler, paymentHandler, os.Getenv("JWT_SECRET"))

	port := os.Getenv("PORT")
	log.Printf("server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
