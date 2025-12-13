package models

type RegisterRequest struct {
    FullName     string `json:"fullname"`
    Email        string `json:"email"`
    Password     string `json:"password"`
    Gender       string `json:"gender"`
    DateOfBirth  string `json:"date_of_birth"`
    Height       *float64 `json:"height"`
    Weight       *float64 `json:"weight"`
    ContactInfo  string `json:"phone_number"`
    Preference   string `json:"preference"`
    HealthGoal   string `json:"health_goal"`
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}