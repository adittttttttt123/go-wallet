package model

import (
	"time"
)

type Wallet struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"uniqueIndex;not null;size:100" json:"user_id"`
	Balance   float64   `gorm:"type:decimal(15,2);default:0" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WalletID  uint      `gorm:"index;not null" json:"wallet_id"`
	Type      string    `gorm:"type:enum('TOP_UP', 'TRANSFER_OUT', 'TRANSFER_IN');not null" json:"type"`
	Amount    float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Reference string    `gorm:"type:varchar(100)" json:"reference"` // e.g., to_user_id or from_user_id
	CreatedAt time.Time `json:"created_at"`
}
