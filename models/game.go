package models

type Street int

const (
	PreFlop Street = iota
	Flop
	Turn
	River
)

type GameState struct {
	ID            int          `json:"id"`
	Players       []Player     `json:"players"`
	Board         []Card       `json:"board"`
	Street        Street       `json:"street"`
	Done          bool         `json:"done"`
	RemainingDeck []Card       `json:"remaining_deck"`
	Winners       []PlayerHand `json:"winners"`
}

type EquityPrediction struct {
	Cards  []Card  `json:"cards"`
	Equity float64 `json:"equity"`
	ID     string  `json:"id"`
}

type HandResult struct {
	Hand        HandRank `json:"winner_hand"`
	TieBreakers []Card   `json:"tiebreakers"`
}

type PlayerHand struct {
	PlayerID string
	Result   HandResult
}

type Cards = JSONColumn[Card]
type Players = JSONColumn[Player]
type PlayerHands = JSONColumn[PlayerHand]
