package engine

import (
	"math/rand/v2"
	"poker-engine/models"
)

func ShuffleDeck() models.Deck {
	deck := models.NewDeck()

	rand.Shuffle(len(deck), func(i int, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

func DealPlayer(remainingDeck models.Deck) (models.Deck, []models.Card) {
	return dealCards(2, remainingDeck)
}

func DealCommunity(remainingDeck models.Deck, street models.Street) (models.Deck, []models.Card) {
	if street == models.Flop {
		return dealCards(3, remainingDeck)
	}

	return dealCards(1, remainingDeck)
}

func dealCards(amount int, deck models.Deck) (models.Deck, []models.Card) {
	cards := make([]models.Card, 0, amount)

	for range amount {
		cards = append(cards, deck[0])
		deck = deck[1:]
	}

	return deck, cards
}
