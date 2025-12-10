package models

import "time"

type Recommendation struct {
    RecommendationID uint       `gorm:"primaryKey" json:"recommendation_id"`
    UserID           uint       `gorm:"not null" json:"user_id"`
    Message          string     `json:"message"`
    Category         string     `json:"category"`
    RiskLevel        *string    `json:"risk_level,omitempty"`
    CreatedAt        time.Time  `json:"created_at"`
}

