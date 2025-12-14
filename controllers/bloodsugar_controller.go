package controllers

import (
	"net/http"
	"sweetake/database"
	"sweetake/models"
	"sweetake/utils"

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
	claimsData, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := claimsData.(*utils.Claims)

	var metrics []models.BloodSugarMetric
	if err := database.DB.Where("user_id = ?", claims.UserID).Order("date_time desc").Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
