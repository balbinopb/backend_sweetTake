package controllers

import (
	"net/http"
	"strconv"
	"time"

	"sweetake/database"
	"sweetake/models"

	"gorm.io/datatypes"

	"github.com/gin-gonic/gin"
)

// parseDate parses YYYY-MM-DD
func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

// GET /charts/daily?user_id=1&start=2025-02-01&end=2025-02-07
func CreateDailyConsumptionChart(c *gin.Context) {
	userIDStr := c.Query("user_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	if userIDStr == "" || startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, start, and end are required query params"})
		return
	}

	uid64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	userID := uint(uid64)

	start, err := parseDate(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date, expected YYYY-MM-DD"})
		return
	}
	end, err := parseDate(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date, expected YYYY-MM-DD"})
		return
	}

	// fetch consumptions in range
	var cons []models.Consumption
	if err := database.DB.
		Where("user_id = ? AND date_time BETWEEN ? AND ?", userID, start, end).
		Order("date_time ASC").
		Find(&cons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load consumptions"})
		return
	}

	// build chart JSON
	data, err := BuildDailyConsumptionChart(cons, userID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build chart data"})
		return
	}

	// save snapshot graph
	graph := models.Graph{
		UserID:     userID,
		GraphType:  "Daily Intake",
		StartDate:  &start,
		EndDate:    &end,
		DataPoints: datatypes.JSON(data),
	}
	if err := database.DB.Create(&graph).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save graph"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "daily chart snapshot created",
		"graph_id":    graph.GraphID,
		"graph_type":  graph.GraphType,
		"user_id":     graph.UserID,
		"start_date":  startStr,
		"end_date":    endStr,
		"data_points": jsonRaw(data),
	})
}

// GET /charts/weekly?user_id=1&start=2025-02-01&end=2025-03-01
func CreateWeeklyConsumptionChart(c *gin.Context) {
	userIDStr := c.Query("user_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	if userIDStr == "" || startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, start, and end are required query params"})
		return
	}

	uid64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	userID := uint(uid64)

	start, err := parseDate(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date, expected YYYY-MM-DD"})
		return
	}
	end, err := parseDate(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date, expected YYYY-MM-DD"})
		return
	}

	// fetch consumptions
	var cons []models.Consumption
	if err := database.DB.
		Where("user_id = ? AND date_time BETWEEN ? AND ?", userID, start, end).
		Order("date_time ASC").
		Find(&cons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load consumptions"})
		return
	}

	// build chart JSON
	data, err := BuildWeeklyConsumptionChart(cons, userID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build chart data"})
		return
	}

	// save snapshot graph
	graph := models.Graph{
		UserID:     userID,
		GraphType:  "Weekly Trend",
		StartDate:  &start,
		EndDate:    &end,
		DataPoints: datatypes.JSON(data),
	}
	if err := database.DB.Create(&graph).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save graph"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "weekly chart snapshot created",
		"graph_id":    graph.GraphID,
		"graph_type":  graph.GraphType,
		"user_id":     graph.UserID,
		"start_date":  startStr,
		"end_date":    endStr,
		"data_points": jsonRaw(data),
	})
}

// GET /charts/monthly?user_id=1&start=2025-01-01&end=2025-12-31
func CreateMonthlyConsumptionChart(c *gin.Context) {
	userIDStr := c.Query("user_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	if userIDStr == "" || startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, start, and end are required query params"})
		return
	}

	uid64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	userID := uint(uid64)

	start, err := parseDate(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date, expected YYYY-MM-DD"})
		return
	}
	end, err := parseDate(endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date, expected YYYY-MM-DD"})
		return
	}

	// fetch consumptions
	var cons []models.Consumption
	if err := database.DB.
		Where("user_id = ? AND date_time BETWEEN ? AND ?", userID, start, end).
		Order("date_time ASC").
		Find(&cons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load consumptions"})
		return
	}

	// build chart JSON
	data, err := BuildMonthlyConsumptionChart(cons, userID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build chart data"})
		return
	}

	// save snapshot graph
	graph := models.Graph{
		UserID:     userID,
		GraphType:  "Monthly Summary",
		StartDate:  &start,
		EndDate:    &end,
		DataPoints: datatypes.JSON(data),
	}
	if err := database.DB.Create(&graph).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save graph"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "monthly chart snapshot created",
		"graph_id":    graph.GraphID,
		"graph_type":  graph.GraphType,
		"user_id":     graph.UserID,
		"start_date":  startStr,
		"end_date":    endStr,
		"data_points": jsonRaw(data),
	})
}

// jsonRaw allows returning the raw JSON array under "data_points" without double-encoding.
func jsonRaw(b []byte) gin.H {
	return gin.H{"raw": string(b)}
}

// package controllers

// import (
// 	"time"

// 	"sweetake/database"
// 	"sweetake/models"

// 	"gorm.io/datatypes"
// )

// // Save a WEEKLY chart for a user & period
// func SaveWeeklyConsumptionGraph(cons []models.Consumption, userID uint, start, end time.Time) error {
// 	data, err := BuildWeeklyConsumptionChart(cons, userID, start, end)
// 	if err != nil {
// 		return err
// 	}

// 	// snapshot graph record
// 	graph := models.Graph{
// 		UserID:     userID,
// 		GraphType:  "Weekly Trend",
// 		StartDate:  &start,
// 		EndDate:    &end,
// 		DataPoints: datatypes.JSON(data), // <-- JSON bytes
// 	}

// 	return database.DB.Create(&graph).Error
// }

// // Similarly for DAILY
// func SaveDailyConsumptionGraph(cons []models.Consumption, userID uint, start, end time.Time) error {
// 	data, err := BuildDailyConsumptionChart(cons, userID, start, end)
// 	if err != nil {
// 		return err
// 	}
// 	graph := models.Graph{
// 		UserID:     userID,
// 		GraphType:  "Daily Intake",
// 		StartDate:  &start,
// 		EndDate:    &end,
// 		DataPoints: datatypes.JSON(data),
// 	}
// 	return database.DB.Create(&graph).Error
// }

// // And MONTHLY
// func SaveMonthlyConsumptionGraph(cons []models.Consumption, userID uint, start, end time.Time) error {
// 	data, err := BuildMonthlyConsumptionChart(cons, userID, start, end)
// 	if err != nil {
// 		return err
// 	}
// 	graph := models.Graph{
// 		UserID:     userID,
// 		GraphType:  "Monthly Summary",
// 		StartDate:  &start,
// 		EndDate:    &end,
// 		DataPoints: datatypes.JSON(data),
// 	}
// 	return database.DB.Create(&graph).Error
// }
