package models

type ConsumptionRequest struct {
	UserID     int     `json:"user_id" binding:"required"`
	FoodType   string  `json:"type" binding:"max=255"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	SugarGrams float64 `json:"sugar_grams" binding:"required,gt=0"`
	Context    string  `json:"context" binding:"required"`
}
