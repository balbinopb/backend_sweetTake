package controllers

import (
	"net/http"

	"sweetake/models"
	"sweetake/database"
	"sweetake/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)




// REGISTER
func Register(c *gin.Context) {
	var input models.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// duplicate check
	var existing models.User
	if err := database.DB.Where("email = ? OR username = ?", input.Email, input.Username).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
		return
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	user := models.User{
			Username:    input.Username,
			Email:       input.Email,
			Password:    string(hashedPassword),
			FullName:    &input.FullName,
			Gender:      input.Gender,
			DateOfBirth: input.DateOfBirth,
			Height:      input.Height,
			Weight:      input.Weight,
			ContactInfo: input.ContactInfo,
		}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})

}

// LOGIN
func Login(c *gin.Context) {
	var input models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// CHECK PASSWORD
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// GENERATE JWT
	token, err := utils.GenerateJWT(user.UserID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"token": token})
}

// PROFILE (Protected)
func Profile(c *gin.Context) {
	claimsData, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}

	claims := claimsData.(*utils.Claims)

	var user models.User
	if err := database.DB.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"user": gin.H{
			"username":      user.Username,
			"email":         user.Email,
			"fullname":      user.FullName,
			"gender":        user.Gender,
			"date_of_birth": user.DateOfBirth,
			"height":        user.Height,
			"weight":        user.Weight,
			"contact_info":  user.ContactInfo,
		},
	})
}