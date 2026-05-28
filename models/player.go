package models

import (
	"github.com/0x6flab/namegenerator"
	"github.com/google/uuid"
)

type Player struct {
	ID         string  `json:"id"`
	HoleCards  [2]Card `json:"holeCards"`
	Folded     bool    `json:"folded"`
	Equity     float64 `json:"equity"`
	Winner     bool    `json:"winner"`
	PlayerName string  `json:"player_name"`
}

func NewPlayer() (*Player, error) {
	newPlayerId, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	nameGenie := namegenerator.NewGenerator()

	return &Player{
		ID:         newPlayerId.String(),
		HoleCards:  [2]Card{},
		Folded:     false,
		Equity:     0.0,
		Winner:     false,
		PlayerName: nameGenie.Generate(),
	}, nil
}
