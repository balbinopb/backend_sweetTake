package database

import (
	"log"
	"sweetake/models"
)

func DBMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Consumption{},
		&models.BloodSugarMetric{},
		&models.Graph{},
		&models.Recommendation{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}
