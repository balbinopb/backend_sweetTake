package controllers

import (
	"errors"
	"net/http"
	"strings"
	"sweetake/models"
	"sweetake/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// memory storage just temporary
var (
	users      = []models.User{}
	userID     = 1
	usersMutex sync.Mutex
)

// --- CONTROLLERS ---

// REGISTER
func Register(c *gin.Context) {
	var input models.RegisterRequest

	// if error
	err := c.ShouldBindJSON(&input);
	if err != nil {

		// handle validation errors
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]string, len(ve))
			for i, fe := range ve {
				field := fe.Field()
				tag := fe.Tag()
				switch field {
				case "Username":
					out[i] = "Username is required"
				case "Email":
					switch tag {
					case "required":
						out[i] = "Email is required"
					case "email":
						out[i] = "Invalid email format"
					}
				case "Password":
					switch tag {
					case "required":
						out[i] = "Password is required"
					case "min":
						out[i] = "Password must be at least 6 characters long"
					}
				default:
					out[i] = fe.Error()
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(out, ", ")})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check duplicates
	for _, u := range users {
		if u.Username == input.Username {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		if u.Email == input.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// Save user
	usersMutex.Lock()
	newUser := models.User{
		ID:       userID,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
	}
	userID++
	users = append(users, newUser)
	usersMutex.Unlock()

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// LOGIN
func Login(c *gin.Context) {
	var input models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var found *models.User
	for _, u := range users {
		if u.Email == input.Email {
			found = &u
			break
		}
	}

	if found == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate JWT
	token, _ := utils.GenerateJWT(found.Email)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// protected to fetched/get
func Profile(c *gin.Context) {
	claimsData, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
		return
	}

	claims, ok := claimsData.(*utils.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims type"})
		return
	}

	email := claims.Email

	// Find user by email
	var found *models.User
	for _, u := range users {
		if u.Email == email {
			found = &u
			break
		}
	}

	if found == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully",
		"user": gin.H{
			"username": found.Username,
			"email":    found.Email,
		},
	})
}
