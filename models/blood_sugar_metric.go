package models

import "time"

type BloodSugarMetric struct {
	MetricID uint `gorm:"primaryKey" json:"metric_id"`
	UserID   uint `gorm:"not null" json:"user_id"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	DateTime        time.Time `json:"date_time"`
	BloodSugarData  float64   `json:"blood_sugar"`
	Context         string    `json:"context"`
	CreatedAt       time.Time `json:"created_at"`
}
