package models

import "time"

type RegisterRequest struct {
	FullName    string   `json:"fullname"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Gender      string   `json:"gender"`
	DateOfBirth string   `json:"date_of_birth"`
	Height      *float64 `json:"height"`
	Weight      *float64 `json:"weight"`
	ContactInfo string   `json:"phone_number"`
	Preference  string   `json:"preference"`
	HealthGoal  string   `json:"health_goal"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ConsumptionRequest struct {
	DateTime  time.Time `json:"date_time"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	SugarData float64   `json:"sugar_data"`
	Context   string    `json:"context"`
}

type BloodSugarRequest struct {
	DateTime       time.Time `json:"date_time"`
	BloodSugarData float64   `json:"blood_sugar"`
	Context        string    `json:"context"`
}




type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}



type UpdateProfileRequest struct {
	FullName    *string  `json:"fullname"`
	Gender      *string  `json:"gender"`
	DateOfBirth *string  `json:"date_of_birth"`
	Height      *float64 `json:"height"`
	Weight      *float64 `json:"weight"`
	Preference  *string  `json:"preference"`
	HealthGoal  *string  `json:"health_goal"`
	ContactInfo *string  `json:"phone_number"`
}