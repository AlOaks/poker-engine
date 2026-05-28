package routes

import "poker-engine/models"

type EquityPredictionResponse struct {
	GameState models.GameState
	Equities  []models.EquityPrediction
}
