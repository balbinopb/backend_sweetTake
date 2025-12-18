package models

import "time"

type User struct {
	UserID uint `gorm:"primaryKey" json:"user_id"`

	FullName *string `json:"full_name,omitempty"`
	Email    string  `gorm:"unique;not null" json:"email"`
	Password string  `json:"-" gorm:"size:255;not null"`

	Gender      string     `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Height      *float64   `json:"height,omitempty"`
	Weight      *float64   `json:"weight,omitempty"`
	ContactInfo string     `json:"contact_info,omitempty"`

	MyPreference string `json:"my_preference,omitempty"`
	MyHealthGoal string `json:"my_health_goal,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ResetToken     *string
	ResetExpiresAt *time.Time
}
