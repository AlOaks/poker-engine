package models

type Street int

const (
	Flop Street = iota
	Turn
	River
)

type GameState struct {
	Players       []Player `json:"players"`
	Board         []Card   `json:"board"`
	Street        Street   `json:"street"`
	Done          bool     `json:"done"`
	RemainingDeck []Card   `json:"remaining_deck"`
}
