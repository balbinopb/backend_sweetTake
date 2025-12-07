package models

import "time"

type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FullName   string `json:"full_name" binding:"required"`
	PersonalID string `json:"personal_id"`

	Gender      string     `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Height      *float64   `json:"height"`
	Weight      *float64   `json:"weight"`
	ContactInfo string     `json:"contact_info"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}