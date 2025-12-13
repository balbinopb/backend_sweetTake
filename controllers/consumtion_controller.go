package controllers

import (
	"net/http"
	"sweetake/database"
	"sweetake/models"

	"github.com/gin-gonic/gin"
)

func ConsumptionForm(c *gin.Context) {
	var input models.Consumption

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	consumption := models.Consumption{
		UserID:    userID.(uint),
		DateTime:  input.DateTime,
		Type:      input.Type,
		Amount:    input.Amount,
		SugarData: input.SugarData,
		Context:   input.Context,
	}

	if err := database.DB.Create(&consumption).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create consumption"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "consumption recorded successfully",
	})
}
