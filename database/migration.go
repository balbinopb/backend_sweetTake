package database

import (
	"log"
	"sweetake/models"
)

func DBMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Consumption{},
		&models.BloodSugarMetric{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}
