package engine

import (
	"math/rand/v2"
	"poker-engine/models"
)

func ShuffleDeck() models.Deck {
	return Shuffle(models.NewDeck())
}

func DealPlayer(remainingDeck models.Deck) (models.Deck, []models.Card) {
	return DealCards(2, remainingDeck)
}

func DealCommunity(remainingDeck models.Deck, street models.Street) (models.Deck, []models.Card) {
	if street == models.PreFlop {
		return DealCards(3, remainingDeck)
	}

	return DealCards(1, remainingDeck)
}

func DealCards(amount int, deck models.Deck) (models.Deck, []models.Card) {
	cards := make([]models.Card, 0, amount)

	for range amount {
		cards = append(cards, deck[0])
		deck = deck[1:]
	}

	return deck, cards
}

func Shuffle(deck models.Deck) models.Deck {
	rand.Shuffle(len(deck), func(i int, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

func NewGame() (*models.GameState, error) {
	newDeck := ShuffleDeck()

	var players []models.Player
	for range 6 {
		newPlayer, err := models.NewPlayer()
		if err != nil {
			return nil, err
		}

		remainingDeck, playerHoleCards := DealPlayer(newDeck)
		newPlayer.HoleCards[0], newPlayer.HoleCards[1] = playerHoleCards[0], playerHoleCards[1]
		newDeck = remainingDeck
		players = append(players, *newPlayer)
	}

	return &models.GameState{
		Players:       players,
		Board:         []models.Card{},
		Street:        models.PreFlop,
		Done:          false,
		RemainingDeck: newDeck,
		Winners:       []models.PlayerHand{},
	}, nil
}

func NextGameStep(gameState models.GameState) *models.GameState {
	switch gameState.Street {
	case models.River:
		var evaluatedHands []models.PlayerHand
		for _, plyr := range gameState.Players {
			plyrCompleteHand := append(plyr.HoleCards[:], gameState.Board...)
			evaluatedHand := EvaluateHand(plyrCompleteHand)

			evaluatedHands = append(evaluatedHands, models.PlayerHand{
				Result:   evaluatedHand,
				PlayerID: plyr.ID,
			})
		}

		winnerHands := CompareAllHands(evaluatedHands)
		gameState.Winners = winnerHands
		gameState.Done = true
	default:
		deck, cards := DealCommunity(gameState.RemainingDeck, gameState.Street)
		gameState.Board = append(gameState.Board, cards...)
		gameState.RemainingDeck = deck
		gameState.Street++
	}

	return &gameState
}
