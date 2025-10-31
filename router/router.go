package router

import (
	"sweetake/controllers"
	"sweetake/middleware"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// Group utama /v1/api
	api := r.Group("/v1/api")

	// Auth routes tanpa middleware
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	// Protected routes dengan middleware JWT token
	auth := api.Group("/auth")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/profile", controllers.Profile)
	}

	return r
}
