package controllers

import (
	"net/http"
	"time"

	"sweetake/database"
	"sweetake/models"
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
	if err := database.DB.Where("email = ? OR username = ?", input.Email, input.FullName).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
		return
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	parsedDOB, err := time.Parse(time.RFC3339, input.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date_of_birth must be RFC3339",
		})
		return
	}

	user := models.User{
		FullName:     &input.FullName,
		Email:        input.Email,
		Password:     string(hashedPassword),
		Gender:       input.Gender,
		DateOfBirth:  &parsedDOB,
		Height:       input.Height,
		Weight:       input.Weight,
		MyPreference: input.Preference,
		MyHealthGoal: input.HealthGoal,
		ContactInfo:  input.ContactInfo,
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
			"fullname":      user.FullName,
			"email":         user.Email,
			"gender":        user.Gender,
			"date_of_birth": user.DateOfBirth,
			"height":        user.Height,
			"weight":        user.Weight,
			"contact_info":  user.ContactInfo,
		},
	})
}

func ForgotPassword(c *gin.Context) {
	var input models.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// SECURITY: do not reveal email existence
		c.JSON(http.StatusOK, gin.H{
			"message": "if the email exists, a reset link has been sent",
		})
		return
	}

	// generate token
	token := utils.GenerateRandomToken(2)
	expiry := time.Now().Add(15 * time.Minute)

	user.ResetToken = &token
	user.ResetExpiresAt = &expiry

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save reset token",
		})
		return
	}

	// SEND EMAIL
	err := utils.SendResetEmail(user.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send reset email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "if the email exists, a reset link has been sent",
	})
}

func ResetPassword(c *gin.Context) {
	var input models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("reset_token = ?", input.Token).
		First(&user).Error; err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	if user.ResetExpiresAt == nil || time.Now().After(*user.ResetExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token expired"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(input.NewPassword),
		bcrypt.DefaultCost,
	)

	user.Password = string(hashedPassword)
	user.ResetToken = nil
	user.ResetExpiresAt = nil

	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "password reset successful",
	})
}
