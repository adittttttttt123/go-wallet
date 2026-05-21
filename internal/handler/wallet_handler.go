package handler

import (
	"net/http"
	"wallet-backend/internal/model"
	"wallet-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	repo *repository.WalletRepository
}

func NewWalletHandler(repo *repository.WalletRepository) *WalletHandler {
	return &WalletHandler{repo: repo}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.repo.CreateWallet(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create wallet or user already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "wallet created successfully",
		"data": gin.H{
			"id":         wallet.ID,
			"user_id":    wallet.UserID,
			"balance":    wallet.Balance,
			"created_at": wallet.CreatedAt,
		},
	})
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID := c.Param("user_id")
	wallet, err := h.repo.GetWalletByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": wallet.UserID,
		"balance": wallet.Balance,
	})
}

func (h *WalletHandler) TopUp(c *gin.Context) {
	var req struct {
		UserID string  `json:"user_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.TopUp(req.UserID, req.Amount); err != nil {
		if err.Error() == "wallet not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "top-up successful", "amount": req.Amount})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	var req struct {
		FromUserID string  `json:"from_user_id" binding:"required"`
		ToUserID   string  `json:"to_user_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Transfer(req.FromUserID, req.ToUserID, req.Amount); err != nil {
		if err.Error() == "sender wallet not found" || err.Error() == "receiver wallet not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "insufficient balance" || err.Error() == "cannot transfer to the same account" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}

func (h *WalletHandler) GetHistory(c *gin.Context) {
	userID := c.Param("user_id")
	transactions, err := h.repo.GetTransactions(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Make sure we return an empty array instead of null if there are no transactions
	if transactions == nil {
		transactions = []model.Transaction{}
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"transactions": transactions,
	})
}
