package api

import (
	"net/http"

	"wallet-backend/config"
	"wallet-backend/core/handler"
	"wallet-backend/core/repository"

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
	r.GET("/", func(c *gin.Context) {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Go-Wallet API Dashboard</title>
			<style>
				body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #0f172a; color: #e2e8f0; margin: 0; padding: 20px; }
				.container { max-width: 800px; margin: 0 auto; background: #1e293b; padding: 30px; border-radius: 12px; box-shadow: 0 10px 15px -3px rgba(0,0,0,0.5); }
				h1 { color: #38bdf8; text-align: center; margin-bottom: 5px; }
				.status { text-align: center; color: #10b981; font-weight: bold; margin-bottom: 30px; }
				.endpoint { background: #334155; margin: 15px 0; padding: 15px; border-radius: 8px; border-left: 5px solid #38bdf8; transition: transform 0.2s; }
				.endpoint:hover { transform: translateX(5px); }
				.method { font-weight: bold; padding: 4px 10px; border-radius: 4px; color: white; display: inline-block; width: 60px; text-align: center; margin-right: 10px; font-size: 0.9em; }
				.get { background-color: #10b981; }
				.post { background-color: #3b82f6; }
				code { background: #0f172a; padding: 4px 8px; border-radius: 4px; font-family: monospace; color: #f8fafc; font-size: 0.9em; }
				.desc { margin-top: 10px; color: #cbd5e1; font-size: 0.95em; line-height: 1.5; }
				.footer { text-align: center; margin-top: 40px; font-size: 0.9em; color: #94a3b8; border-top: 1px solid #334155; padding-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>🚀 Go-Wallet API</h1>
				<div class="status">● API is active and running smoothly</div>
				
				<h2>Available Endpoints</h2>
				
				<div class="endpoint">
					<div><span class="method post">POST</span> <code>/api/v1/wallets</code></div>
					<div class="desc">Create a new wallet. Requires JSON body:<br><code>{"user_id": "string"}</code></div>
				</div>
				
				<div class="endpoint">
					<div><span class="method get">GET</span> <code>/api/v1/wallets/:user_id/balance</code></div>
					<div class="desc">Get wallet balance for a specific user.</div>
				</div>
				
				<div class="endpoint">
					<div><span class="method post">POST</span> <code>/api/v1/wallets/topup</code></div>
					<div class="desc">Top up a wallet. Requires JSON body:<br><code>{"user_id": "string", "amount": number}</code></div>
				</div>
				
				<div class="endpoint">
					<div><span class="method post">POST</span> <code>/api/v1/wallets/transfer</code></div>
					<div class="desc">Transfer funds between wallets. Requires JSON body:<br><code>{"from_user_id": "string", "to_user_id": "string", "amount": number}</code></div>
				</div>
				
				<div class="endpoint">
					<div><span class="method get">GET</span> <code>/api/v1/wallets/:user_id/history</code></div>
					<div class="desc">Get transaction history for a specific user.</div>
				</div>
				
				<div class="footer">
					<p>To test these endpoints, import the <code>go-wallet.postman_collection.json</code> file into your <strong>Postman</strong> app.</p>
				</div>
			</div>
		</body>
		</html>
		`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

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
