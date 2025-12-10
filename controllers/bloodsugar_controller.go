package controllers

import (
	"net/http"
	"sweetake/database"
	"sweetake/models"
	"sweetake/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CREATE BLOOD SUGAR METRIC
func CreateBloodSugarMetric(c *gin.Context) {
	var input models.BloodSugarMetric

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from JWT claims
	claimsData, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := claimsData.(*utils.Claims)
	input.UserID = claims.UserID
	input.CreatedAt = time.Now()

	// Validate user exists
	var user models.User
	if err := database.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
		return
	}

	// Parse MeasureDate + MeasureTime into DateTime
	dateParts := strings.Split(input.MeasureDate, "-")
	if len(dateParts) != 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid measure_date format"})
		return
	}
	year, _ := strconv.Atoi(dateParts[0])
	month, _ := strconv.Atoi(dateParts[1])
	day, _ := strconv.Atoi(dateParts[2])

	timeParts := strings.Split(input.MeasureTime, ":")
	if len(timeParts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid measure_time format"})
		return
	}
	hour, _ := strconv.Atoi(timeParts[0])
	minute, _ := strconv.Atoi(timeParts[1])

	input.DateTime = time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)

	// Save to database
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "blood sugar metric created",
		"metric_id": input.MetricID,
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
