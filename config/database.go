package config

import (
	"fmt"
	"log"

	"wallet-backend/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// Menggunakan SQLite agar tidak perlu XAMPP/MySQL lagi
	db, err := gorm.Open(sqlite.Open("wallet.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successful (using SQLite)")

	// Auto Migrate
	err = db.AutoMigrate(&model.Wallet{}, &model.Transaction{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	fmt.Println("Database migration completed")

	return db
}
