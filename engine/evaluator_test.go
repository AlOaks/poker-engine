package engine

import (
	"math/rand/v2"
	"poker-engine/models"
	"testing"
)

func TestEvaluateRoyalFlush(t *testing.T) {
	royalFlushCards := make([]models.Card, 0, 7)

	for i := 14; i >= 8; i-- {
		royalFlushCards = append(royalFlushCards, models.Card{Rank: models.Rank(i), Suit: models.Diamonds})
	}

	evaluation := EvaluateHand(royalFlushCards)

	if evaluation.Hand != models.RoyalFlush {
		t.Fatalf("Wrong evaluation: %s", evaluation.Hand.String())
	}
}

func TestEvaluateStraightFlush(t *testing.T) {
	straightFlushCards := make([]models.Card, 0, 7)

	for i := 13; i >= 7; i-- {
		if len(straightFlushCards) >= 5 {
			straightFlushCards = append(straightFlushCards, models.Card{Rank: models.Rank(i - 2), Suit: models.Spades})

		} else {
			straightFlushCards = append(straightFlushCards, models.Card{Rank: models.Rank(i), Suit: models.Diamonds})

		}
	}

	evaluation := EvaluateHand(straightFlushCards)

	if evaluation.Hand != models.StraightFlush {
		t.Fatalf("Wrong evaluation: %s", evaluation.Hand.String())
	}
}

func TestEvaluateFourOfKind(t *testing.T) {
	fourOfKindCards := make([]models.Card, 0, 7)

	for i := 0; i < 7; i++ {
		if len(fourOfKindCards) <= 3 {
			fourOfKindCards = append(fourOfKindCards, models.Card{Rank: models.Ace, Suit: models.Suit(i)})
		} else {
			fourOfKindCards = append(fourOfKindCards, models.Card{Rank: models.Rank(i), Suit: models.Diamonds})
		}
	}

	evaluation := EvaluateHand(fourOfKindCards)

	if evaluation.Hand != models.FourOfKind {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.FourOfKind.String())
	}
}

func TestEvaluateFullHouse(t *testing.T) {
	fullHouseCards := make([]models.Card, 0, 7)

	// Get TOAK
	for i := 0; i < 3; i++ {
		fullHouseCards = append(fullHouseCards, models.Card{Rank: models.Ace, Suit: models.Suit(i)})
	}

	// Get Pair
	for i := 0; i < 2; i++ {
		fullHouseCards = append(fullHouseCards, models.Card{Rank: models.Two, Suit: models.Suit(i)})
	}

	fullHouseCards = append(fullHouseCards, models.Card{Rank: models.Four, Suit: models.Hearts})
	fullHouseCards = append(fullHouseCards, models.Card{Rank: models.Six, Suit: models.Hearts})

	evaluation := EvaluateHand(fullHouseCards)

	if evaluation.Hand != models.FullHouse {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.FullHouse.String())
	}
}

func TestEvaluateFlush(t *testing.T) {
	flushCards := make([]models.Card, 0, 7)

	// Get flush
	for i := 0; i < 5; i++ {
		randomRank := models.Rank(rand.IntN(12) + 2)
		flushCards = append(flushCards, models.Card{Rank: randomRank, Suit: models.Clubs})
	}

	// Get random
	for i := 0; i < 2; i++ {
		flushCards = append(flushCards, models.Card{Rank: models.Two, Suit: models.Suit(i)})
	}

	evaluation := EvaluateHand(flushCards)

	if evaluation.Hand != models.Flush {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.Flush.String())
	}
}

func TestEvaluateStraight(t *testing.T) {
	straightCards := make([]models.Card, 0, 7)

	for i := 0; i < 7; i++ {
		if len(straightCards) < 5 {
			randomSuit := models.Suit(rand.IntN(3))
			straightCards = append(straightCards, models.Card{Rank: models.Rank(i + 2), Suit: randomSuit})
		} else {
			straightCards = append(straightCards, models.Card{Rank: models.Rank(i), Suit: models.Suit(i - 2)})
		}
	}

	evaluation := EvaluateHand(straightCards)

	if evaluation.Hand != models.Straight {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.Straight)
	}
}

func TestEvaluateThreeOfKind(t *testing.T) {
	threeOfKindCards := make([]models.Card, 0, 7)

	for i := 0; i < 7; i++ {
		if len(threeOfKindCards) < 3 {
			threeOfKindCards = append(threeOfKindCards, models.Card{Rank: models.Ace, Suit: models.Suit(i)})
		} else {
			threeOfKindCards = append(threeOfKindCards, models.Card{Rank: models.Rank(i), Suit: models.Suit(i - 3)})
		}
	}

	evaluation := EvaluateHand(threeOfKindCards)

	if evaluation.Hand != models.ThreeOfKind {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.ThreeOfKind.String())
	}
}

func TestEvaluateTwoPairs(t *testing.T) {
	twoPairsCards := make([]models.Card, 0, 7)

	for i := 0; i < 4; i++ {
		if len(twoPairsCards) >= 2 {
			twoPairsCards = append(twoPairsCards, models.Card{Rank: models.Ten, Suit: models.Suit(i)})
		} else {
			twoPairsCards = append(twoPairsCards, models.Card{Rank: models.Jack, Suit: models.Suit(i)})
		}
	}

	// Get random
	for i := 0; i < 3; i++ {
		twoPairsCards = append(twoPairsCards, models.Card{Rank: models.Rank(i + 2), Suit: models.Suit(i)})
	}

	evaluation := EvaluateHand(twoPairsCards)

	if evaluation.Hand != models.TwoPairs {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.TwoPairs.String())
	}
}

func TestEvaluatePair(t *testing.T) {
	pairCards := make([]models.Card, 0, 7)

	for i := 0; i < 2; i++ {
		pairCards = append(pairCards, models.Card{Rank: models.Jack, Suit: models.Suit(i)})
	}

	// Get random
	for i := 0; i < 5; i++ {
		randomRank := models.Rank(rand.IntN(8) + 2)
		pairCards = append(pairCards, models.Card{Rank: randomRank, Suit: models.Suit(i)})
	}

	evaluation := EvaluateHand(pairCards)

	if evaluation.Hand != models.Pair {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.Pair.String())
	}
}

func TestHighCard(t *testing.T) {
	highCards := make([]models.Card, 0, 7)

	// Get random
	for i := 0; i < 7; i++ {
		randomRank := models.Rank(rand.IntN(8) + 2)
		highCards = append(highCards, models.Card{Rank: randomRank, Suit: models.Suit(i)})
	}

	evaluation := EvaluateHand(highCards)

	t.Logf("%v", highCards)

	if evaluation.Hand != models.HighCard {
		t.Fatalf("Wrong evaluation: %s - Should be: %s", evaluation.Hand.String(), models.HighCard.String())
	}
}
