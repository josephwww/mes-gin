package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	// Use your actual configuration values here
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		AppConfig.Database.Host,
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.Name,
		AppConfig.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
}
