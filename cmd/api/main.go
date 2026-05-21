package main

import (
	"log"
	"net/http"

	"wallet-backend/config"
	"wallet-backend/core/handler"
	"wallet-backend/core/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if exists (ignoring error if it doesn't)
	_ = godotenv.Load()

	// Init DB connection
	db := config.ConnectDatabase()

	// Init Repo & Handler
	walletRepo := repository.NewWalletRepository(db)
	walletHandler := handler.NewWalletHandler(walletRepo)

	// Setup Gin Router
	r := gin.Default()

	// Ping route for health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/wallets", walletHandler.CreateWallet)
		api.GET("/wallets/:user_id/balance", walletHandler.GetBalance)
		api.POST("/wallets/topup", walletHandler.TopUp)
		api.POST("/wallets/transfer", walletHandler.Transfer)
		api.GET("/wallets/:user_id/history", walletHandler.GetHistory)
	}

	// Start server on port 8080
	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
