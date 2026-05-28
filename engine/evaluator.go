package engine

import (
	"poker-engine/models"
	"slices"
	"sort"
)

func EvaluateHand(cards []models.Card) models.HandResult {
	var handResult models.HandResult

	rankCounts := make(map[models.Rank]int)
	suitCounts := make(map[models.Suit]int)

	for _, card := range cards {
		rankCounts[card.Rank]++
		suitCounts[card.Suit]++
	}

	isFlush := false

	// First we arrange the cards from highest to lowest
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank > cards[j].Rank
	})

	// First check to see if its a flush.
	for key, val := range suitCounts {
		if val >= 5 {
			isFlush = true
			flushType, flushCards := processFlush(cards, key)
			handResult.Hand = flushType
			handResult.TieBreakers = flushCards

			break
		}
	}

	if isFlush {
		return handResult
	}

	// Then, we check if its s four-of-a-kind
	for key, val := range rankCounts {
		// If the count is 4, then we know is a four of a kind,
		// because there cannot be more than 4
		if val == 4 {
			handResult.Hand = models.FourOfKind
			handResult.TieBreakers = processFourOfKind(cards, key)

			break
		}
	}

	if handResult.Hand == models.FourOfKind {
		return handResult
	}

	// Ce check if there is a straight and early return if so.
	straightCards, isStraight := isStraight(cards)

	if isStraight {
		return models.HandResult{
			Hand:        models.Straight,
			TieBreakers: straightCards,
		}
	}

	handResult = processSameKind(cards, rankCounts)

	// If it reaches this point, then is just a high card hand
	return handResult

}

func processSameKind(cards []models.Card, rankCounts map[models.Rank]int) models.HandResult {
	var toakGroups []models.Rank
	var pairGroups []models.Rank

	for rank, count := range rankCounts {
		if count == 3 {
			toakGroups = append(toakGroups, rank)
		}

		if count == 2 {
			pairGroups = append(pairGroups, rank)
		}
	}

	// No TOAKs or Pairs
	if len(toakGroups) == 0 && len(pairGroups) == 0 {
		return models.HandResult{
			Hand:        models.HighCard,
			TieBreakers: cards[:5],
		}
	}

	var handResult models.HandResult

	var toaks []models.Card
	var pairs []models.Card
	var kickers []models.Card

	for _, card := range cards {
		if slices.Contains(toakGroups, card.Rank) {
			toaks = append(toaks, card)
		} else if slices.Contains(pairGroups, card.Rank) {
			pairs = append(pairs, card)
		} else {
			kickers = append(kickers, card)
		}
	}

	switch len(kickers) {
	case 0, 2:
		// This case means there are 1 Toak and 1 Pair OR 1 toak and 2 Pairs - thus a FullHouse
		handResult.TieBreakers = append(toaks, pairs[:2]...)
		handResult.Hand = models.FullHouse
	case 1:
		// This case means there are 2 Toaks or 3 pairs - thus a FullHouse or Two pairs
		// Need to determine which is it.
		if len(toaks) == 6 {
			handResult.TieBreakers = toaks[:5]
			handResult.Hand = models.FullHouse
		} else {
			var fifthCard models.Card
			if pairs[4].Rank > kickers[0].Rank {
				fifthCard = pairs[4]
			} else {
				fifthCard = kickers[0]
			}

			handResult.TieBreakers = append(pairs[:4], fifthCard)
			handResult.Hand = models.TwoPairs
		}
	case 3:
		// This case means there are 2 Pairs - thus a TwoPair
		handResult.Hand = models.TwoPairs
		handResult.TieBreakers = append(pairs, kickers[0])
	case 4:
		// This case means there is 1 Toak - thus a ThreeOfKind
		handResult.Hand = models.ThreeOfKind
		handResult.TieBreakers = append(toaks, kickers[:2]...)
	case 5:
		// This case means there is 1 pair - thus a Pair
		handResult.Hand = models.Pair
		handResult.TieBreakers = append(pairs, kickers[:3]...)
	}

	return handResult
}

func processFourOfKind(cards []models.Card, rank models.Rank) []models.Card {
	outCards := make([]models.Card, 0, 2)

	// If the first card is of the rank, then we know the next 3 are the same,
	// therefore we use the 5th card as the kicker
	if cards[0].Rank == rank {
		return cards[:5]
	}

	for _, card := range cards {
		if card.Rank == rank {
			outCards = append(outCards, card)
		}
	}
	return append(outCards, cards[0])
}

func processFlush(cards []models.Card, suit models.Suit) (models.HandRank, []models.Card) {
	flushType := models.Flush
	var outCards []models.Card

	for _, card := range cards {
		if card.Suit == suit {
			outCards = append(outCards, card)
		}
	}

	straightCards, isStraight := isStraight(outCards)

	if isStraight {
		outCards = straightCards
		if outCards[4].Rank == models.Ten {
			flushType = models.RoyalFlush
		} else {
			flushType = models.StraightFlush
		}
	} else {
		flushType = models.Flush
		outCards = outCards[:5]
	}

	return flushType, outCards[:5]
}

func isStraight(cards []models.Card) ([]models.Card, bool) {
	isStraight := false
	var outCards []models.Card

	uniqueCards := models.UniqueSet(cards)

	for i := 0; i <= 2; i++ {
		if i+5 > len(uniqueCards) {
			break
		}

		window := uniqueCards[i : i+5]
		isWindowStraight := true

		for j, card := range window {
			if j+1 < len(window) {
				switch {
				case card.Rank == models.Ace:
					if !slices.Contains([]models.Rank{models.Five, models.King}, window[j+1].Rank) {
						isWindowStraight = false
						break
					}
				default:
					if card.Rank-1 != window[j+1].Rank {
						isWindowStraight = false
						break
					}
				}
			}

		}

		if isWindowStraight {
			outCards = window
			isStraight = true
			break
		}
	}

	return outCards, isStraight
}

// Compare two hands
func CompareTwoHands(handOne, handTwo models.PlayerHand) *models.PlayerHand {
	var winnerHand *models.PlayerHand

	if handOne.Result.Hand == handTwo.Result.Hand {
		// Let's check the kickers
		for i, kicker := range handOne.Result.TieBreakers {
			if kicker.Rank > handTwo.Result.TieBreakers[i].Rank {
				winnerHand = &handOne
				break
			}

			if kicker.Rank < handTwo.Result.TieBreakers[i].Rank {
				winnerHand = &handTwo
				break
			}
		}

	} else if handOne.Result.Hand > handTwo.Result.Hand {
		winnerHand = &handOne
	} else {
		winnerHand = &handTwo
	}

	return winnerHand
}

func CompareAllHands(hands []models.PlayerHand) []models.PlayerHand {
	var winnerHands []models.PlayerHand

	for i, hand := range hands {
		if i == 0 {
			winnerHands = append(winnerHands, hand)
			continue
		}

		currentWinner := winnerHands[0]
		if hand.Result.Hand == currentWinner.Result.Hand {
			winner := CompareTwoHands(hand, currentWinner)
			if winner == nil {
				winnerHands = append(winnerHands, hand)
				continue
			}

			winnerHands = []models.PlayerHand{*winner}
			continue
		}

		if hand.Result.Hand > currentWinner.Result.Hand {
			winnerHands = []models.PlayerHand{hand}
		}
	}

	return winnerHands
}
