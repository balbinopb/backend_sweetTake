package models

import "time"

type Consumption struct {
    ConsumptionID uint      `gorm:"primaryKey" json:"consumption_id"`
    UserID        uint      `gorm:"not null" json:"user_id"`
    DateTime      time.Time `json:"date_time"`
    Type          string    `json:"type"`
    Amount        *float64  `json:"amount"`
    SugarData     *float64  `json:"sugar_data"`
    Context       string    `json:"context"`
    CreatedAt     time.Time `json:"created_at"`

    Graphs []Graph `gorm:"foreignKey:ConsumptionID"`
}
