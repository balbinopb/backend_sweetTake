package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	databaseURL := os.Getenv("DATABASE_URL")

	// Local development 
	if databaseURL == "" {
		log.Println("DATABASE_URL not set, using local database")
		databaseURL = "postgresql://postgres:postgres@localhost:5432/sweettake_db?sslmode=disable"
	} else {
		log.Println("Using production database")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connected successfully")
	return db
}
