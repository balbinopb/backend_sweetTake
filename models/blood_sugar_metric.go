package models

import "time"

type BloodSugarMetric struct {
	MetricID uint `gorm:"primaryKey" json:"metric_id"`
	UserID   uint `gorm:"not null" json:"user_id"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	MeasureDate string    `json:"measure_date"`
	MeasureTime string    `json:"measure_time"`
	DateTime    time.Time `json:"date_time"`
	Value       float64   `json:"value"`
	Context     string    `json:"context"`
	CreatedAt   time.Time `json:"created_at"`
}
