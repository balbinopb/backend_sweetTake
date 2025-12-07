package models

import "time"

type Recommendation struct {
    RecommendationID uint           `gorm:"primaryKey" json:"recommendation_id"`
    UserID           uint           `gorm:"not null" json:"user_id"`
    Message          string         `json:"message"`
    Category         string         `json:"category"`
    RiskLevel        RiskLevelType  `gorm:"type:risk_level_type" json:"risk_level"`
    CreatedAt        time.Time      `json:"created_at"`
}
