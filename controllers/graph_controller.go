package controllers

import (
	"time"

	"sweetake/database"
	"sweetake/models"

	"gorm.io/datatypes"
)

// Save a WEEKLY chart for a user & period
func SaveWeeklyConsumptionGraph(cons []models.Consumption, userID uint, start, end time.Time) error {
	data, err := BuildWeeklyConsumptionChart(cons, userID, start, end)
	if err != nil {
		return err
	}

	// snapshot graph record
	graph := models.Graph{
		UserID:     userID,
		GraphType:  "Weekly Trend",
		StartDate:  &start,
		EndDate:    &end,
		DataPoints: datatypes.JSON(data), // <-- JSON bytes
	}

	return database.DB.Create(&graph).Error
}

// Similarly for DAILY
func SaveDailyConsumptionGraph(cons []models.Consumption, userID uint, start, end time.Time) error {
	data, err := BuildDailyConsumptionChart(cons, userID, start, end)
	if err != nil {
		return err
	}
	graph := models.Graph{
		UserID:     userID,
		GraphType:  "Daily Intake",
		StartDate:  &start,
		EndDate:    &end,
		DataPoints: datatypes.JSON(data),
	}
	return database.DB.Create(&graph).Error
}

// And MONTHLY
func SaveMonthlyConsumptionGraph(cons []models.Consumption, userID uint, start, end time.Time) error {
	data, err := BuildMonthlyConsumptionChart(cons, userID, start, end)
	if err != nil {
		return err
	}
	graph := models.Graph{
		UserID:     userID,
		GraphType:  "Monthly Summary",
		StartDate:  &start,
		EndDate:    &end,
		DataPoints: datatypes.JSON(data),
	}
	return database.DB.Create(&graph).Error
}
