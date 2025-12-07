package models

import (
    "time"
    "gorm.io/datatypes"
)

type Graph struct {
    GraphID        uint           `gorm:"primaryKey" json:"graph_id"`
    ConsumptionID  uint           `gorm:"not null" json:"consumption_id"`
    GraphType      string         `json:"graph_type"`
    DataPoints     datatypes.JSON `json:"data_points"`
    CreatedAt      time.Time      `json:"created_at"`
}
