package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	// Railway provides DATABASE_URL (BEST)
	databaseURL := os.Getenv("DATABASE_URL")

	var dsn string

	if databaseURL != "" {
		// Railway / production
		dsn = databaseURL
	} else {
		// Local development fallback
		host := os.Getenv("PGHOST")
		if host == "" {
			host = "localhost"
		}

		user := os.Getenv("PGUSER")
		if user == "" {
			user = "postgres"
		}

		password := os.Getenv("PGPASSWORD")
		if password == "" {
			password = "postgres"
		}

		dbname := os.Getenv("PGDATABASE")
		if dbname == "" {
			dbname = "sweettake_db"
		}

		port := os.Getenv("PGPORT")
		if port == "" {
			port = "5432"
		}

		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connected successfully")

	return db
}
