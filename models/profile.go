package models

import "time"

type Profile struct {
    ProfileID   uint       `gorm:"primaryKey" json:"profile_id"`
    UserID      uint       `gorm:"unique;not null" json:"user_id"`
    Preferences *string    `json:"preferences,omitempty"` 
    HealthGoals *string    `json:"health_goals,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

