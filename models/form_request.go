package models

type ConsumptionRequest struct {
	UserID     int     `json:"user_id" binding:"required"`
	Type       string  `json:"type" binding:"max=255"`              // Optional free text, max 255 chars
	Amount     float64 `json:"amount" binding:"required,gt=0"`      // Required, must be > 0
	SugarGrams float64 `json:"sugar_grams" binding:"required,gt=0"` // Required, must be > 0
	Context    string  `json:"context" binding:"required"`          // Required (dropdown)
}
