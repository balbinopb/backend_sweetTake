package controllers

import (
	"net/http"
	"sweetake/database"
	"sweetake/models"

	"github.com/gin-gonic/gin"
)

// CREATE BLOOD SUGAR METRIC
func CreateBloodSugarMetric(c *gin.Context) {
	var input models.BloodSugarRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from JWT claims
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bloodSugar := models.BloodSugarMetric{
		UserID:         userID.(uint),
		DateTime:       input.DateTime,
		BloodSugarData: input.BloodSugarData,
		Context:        input.Context,
	}

	if err := database.DB.Create(&bloodSugar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create blood sugar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "blood sugar recorded successfully",
	})

}

// GET SINGLE BLOOD SUGAR METRIC BY ID
func GetBloodSugarMetric(c *gin.Context) {
	id := c.Param("id")
	var metric models.BloodSugarMetric

	if err := database.DB.First(&metric, "metric_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "metric not found"})
		return
	}

	c.JSON(http.StatusOK, metric)
}

// GET ALL BLOOD SUGAR METRICS FOR LOGGED-IN USER
func GetAllBloodSugarMetrics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var bloodSugar []models.BloodSugarMetric

	if err := database.DB.
		Where("user_id = ?", userID).
		Order("date_time DESC").
		Find(&bloodSugar).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch blood sugar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bloodSugar,
	})
}

// DELETE BLOOD SUGAR METRIC
func DeleteBloodSugarMetric(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")

	var metric models.BloodSugarMetric

	// Ensure metric belongs to the user
	if err := database.DB.
		Where("metric_id = ? AND user_id = ?", id, userID).
		First(&metric).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "metric not found",
		})
		return
	}

	if err := database.DB.Delete(&metric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete blood sugar metric",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "blood sugar metric deleted successfully",
	})
}
