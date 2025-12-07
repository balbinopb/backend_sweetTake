package models

import "time"

type BloodSugarMetric struct {
    MetricID     uint      `gorm:"primaryKey" json:"metric_id"`
    UserID       uint      `gorm:"not null" json:"user_id"`
    MeasureDate  time.Time `json:"measure_date"`
    MeasureTime  time.Time `json:"measure_time"`
    Value        float64   `json:"value"`
    Context      string    `json:"context"`
    CreatedAt    time.Time `json:"created_at"`
}
