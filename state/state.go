package state

import "poker-engine/models"

type AppState struct {
	Games map[string]models.GameState
}

func NewAppState() AppState {
	return AppState{Games: make(map[string]models.GameState)}
}
