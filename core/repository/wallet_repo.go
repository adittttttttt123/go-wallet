package repository

import (
	"errors"
	"wallet-backend/core/model"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// 1. Pembuatan Akun & Dompet
func (r *WalletRepository) CreateWallet(userID string) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID:  userID,
		Balance: 0,
	}
	err := r.db.Create(wallet).Error
	return wallet, err
}

func (r *WalletRepository) GetWalletByUserID(userID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

// 2. Proses Top-Up (Isi Saldo)
func (r *WalletRepository) TopUp(userID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet
		// Lock the row for update to prevent race conditions (Pessimistic Locking)
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return errors.New("wallet not found")
		}

		wallet.Balance += amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		transaction := model.Transaction{
			WalletID: wallet.ID,
			Type:     "TOP_UP",
			Amount:   amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

// 3. Proses Transfer (Kirim Uang)
func (r *WalletRepository) Transfer(fromUserID, toUserID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if fromUserID == toUserID {
		return errors.New("cannot transfer to the same account")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var fromWallet model.Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", fromUserID).First(&fromWallet).Error; err != nil {
			return errors.New("sender wallet not found")
		}

		if fromWallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		var toWallet model.Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", toUserID).First(&toWallet).Error; err != nil {
			return errors.New("receiver wallet not found")
		}

		// Deduct from sender
		fromWallet.Balance -= amount
		if err := tx.Save(&fromWallet).Error; err != nil {
			return err
		}

		// Add to receiver
		toWallet.Balance += amount
		if err := tx.Save(&toWallet).Error; err != nil {
			return err
		}

		// Record sender transaction (mutasi keluar)
		txOut := model.Transaction{
			WalletID:  fromWallet.ID,
			Type:      "TRANSFER_OUT",
			Amount:    amount,
			Reference: toUserID,
		}
		if err := tx.Create(&txOut).Error; err != nil {
			return err
		}

		// Record receiver transaction (mutasi masuk)
		txIn := model.Transaction{
			WalletID:  toWallet.ID,
			Type:      "TRANSFER_IN",
			Amount:    amount,
			Reference: fromUserID,
		}
		if err := tx.Create(&txIn).Error; err != nil {
			return err
		}

		return nil
	})
}

// 4. Pencatatan Mutasi (History)
func (r *WalletRepository) GetTransactions(userID string) ([]model.Transaction, error) {
	var wallet model.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, errors.New("wallet not found")
	}

	var transactions []model.Transaction
	err := r.db.Where("wallet_id = ?", wallet.ID).Order("created_at desc").Find(&transactions).Error
	return transactions, err
}
