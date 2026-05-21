package config

import (
	"fmt"
	"log"
	"os"

	"wallet-backend/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// Akan mengambil URL dari Environment Variable (Supabase)
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Println("Peringatan: Variabel DB_DSN kosong!")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successful (PostgreSQL/Supabase)")

	// Auto Migrate
	err = db.AutoMigrate(&model.Wallet{}, &model.Transaction{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	fmt.Println("Database migration completed")

	return db
}
