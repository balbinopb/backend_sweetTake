package models

import "time"

type Profile struct {
    ProfileID    uint           `gorm:"primaryKey" json:"profile_id"`
    UserID       uint           `gorm:"unique;not null" json:"user_id"`
    Preferences  PreferenceType `gorm:"type:preference_type" json:"preferences"`
    HealthGoals  HealthGoalType `gorm:"type:health_goal_type" json:"health_goals"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}
