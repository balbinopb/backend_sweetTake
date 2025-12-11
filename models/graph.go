package models

import (
	"time"

	"gorm.io/datatypes"
)

type Graph struct {
	GraphID uint `gorm:"primaryKey" json:"graph_id"`
	// Make this nullable for non-entry charts:
	ConsumptionID *uint `json:"consumption_id,omitempty"`

	// Add owner & range:
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	GraphType string     `json:"graph_type"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`

	DataPoints datatypes.JSON `json:"data_points"`
	CreatedAt  time.Time      `json:"created_at"`
}
