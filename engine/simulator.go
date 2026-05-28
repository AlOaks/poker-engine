package engine

import (
	"poker-engine/models"
	"slices"
	"sync"
)

type Simulator struct {
	events chan *models.EquityPrediction
}

func NewSimulator() *Simulator {
	channel := make(chan *models.EquityPrediction, 6)

	return &Simulator{
		events: channel,
	}
}

func (s *Simulator) GetChannel() <-chan *models.EquityPrediction {
	return s.events
}

func (s *Simulator) RunPredictions(game models.GameState) {
	var wg sync.WaitGroup

	for _, player := range game.Players {
		if !player.Folded {
			wg.Add(1)
			go func(player models.Player, game models.GameState) {
				defer wg.Done()
				// Calculations
				playerEquity := s.CalculateEquity(player, game)

				s.events <- &playerEquity
			}(player, game)
		}
	}

	go func() {
		wg.Wait()
		close(s.events)
	}()
}

func (s *Simulator) CalculateEquity(player models.Player, game models.GameState) models.EquityPrediction {
	wins := 0.0

	for range 10000 {
		clonedGameState := models.GameState{
			Players:       slices.Clone(game.Players),
			Board:         slices.Clone(game.Board),
			Street:        game.Street,
			Done:          game.Done,
			RemainingDeck: slices.Clone(game.RemainingDeck),
		}

		shuffledDeck := Shuffle(clonedGameState.RemainingDeck)

		switch len(clonedGameState.Board) {
		case 0:
			deck, cards := DealCards(5, shuffledDeck)
			clonedGameState.Board = cards
			clonedGameState.RemainingDeck = deck
		case 3:
			deck, cards := DealCards(2, shuffledDeck)
			clonedGameState.Board = append(clonedGameState.Board, cards...)
			clonedGameState.RemainingDeck = deck
		case 4:
			deck, cards := DealCards(1, shuffledDeck)
			clonedGameState.Board = append(clonedGameState.Board, cards...)
			clonedGameState.RemainingDeck = deck
		}

		var evaluatedHands []models.PlayerHand
		for _, plyr := range clonedGameState.Players {
			if !plyr.Folded {
				plyrCompleteHand := append(plyr.HoleCards[:], clonedGameState.Board...)
				evaluatedHand := EvaluateHand(plyrCompleteHand)

				evaluatedHands = append(evaluatedHands, models.PlayerHand{
					Result:   evaluatedHand,
					PlayerID: plyr.ID,
				})
			}
		}

		winnerHands := CompareAllHands(evaluatedHands)
		for _, winner := range winnerHands {
			if winner.PlayerID == player.ID {
				if len(winnerHands) > 1 {
					wins = wins + 0.5
				} else {
					wins++
				}
			}

		}
	}

	// Calculate equity
	equity := (wins / 10000.0) * 100

	return models.EquityPrediction{
		Cards:  player.HoleCards[:],
		Equity: equity,
		ID:     player.ID,
	}
}
