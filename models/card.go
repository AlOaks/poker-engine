package models

type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Card struct {
	Rank Rank `json:"rank"`
	Suit Suit `json:"suit"`
}

type Deck []Card

func NewDeck() Deck {
	deck := make(Deck, 0, 52)
	for suit := Spades; suit <= Clubs; suit++ {
		for rank := Two; rank <= Ace; rank++ {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}
	return deck
}
