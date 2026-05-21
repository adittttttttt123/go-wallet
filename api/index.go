package api

import (
	"net/http"

	"wallet-backend/config"
	"wallet-backend/internal/handler"
	"wallet-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	// Initialize Gin
	r := gin.Default()

	// Initialize DB (It will panic if DB_DSN is not set properly, which is good for Vercel logs)
	db := config.ConnectDatabase()
	walletRepo := repository.NewWalletRepository(db)
	walletHandler := handler.NewWalletHandler(walletRepo)

	// Setup routes
	api := r.Group("/api/v1")
	{
		api.POST("/wallets", walletHandler.CreateWallet)
		api.GET("/wallets/:user_id/balance", walletHandler.GetBalance)
		api.POST("/wallets/topup", walletHandler.TopUp)
		api.POST("/wallets/transfer", walletHandler.Transfer)
		api.GET("/wallets/:user_id/history", walletHandler.GetHistory)
	}

	app = r
}

// Handler is the main entry point for Vercel Serverless Functions
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
