package engine

import (
	"cmp"
	"math/rand/v2"
	"poker-engine/models"
	"slices"
)

func Folder(game *models.GameState) (*models.GameState, error) {
	switch game.Street {
	case models.PreFlop:
		for i, player := range game.Players {
			if player.Equity < 5.0 {
				game.Players[i].Folded = true
			}
		}
	default:
		if game.Street != models.River {
			foldingCount := rand.IntN(2) + 1
			slices.SortFunc(game.Players, func(a, b models.Player) int {
				return cmp.Compare(a.Equity, b.Equity)
			})

			nonFoldedCount := 0
			for _, player := range game.Players {
				if !player.Folded {
					nonFoldedCount++
				}
			}

			if nonFoldedCount-foldingCount > 2 {
				for i, player := range game.Players {
					if !player.Folded {
						if foldingCount == 0 {
							break
						}
						game.Players[i].Folded = true
						nonFoldedCount--
						foldingCount--
						if nonFoldedCount == 2 {
							break
						}
					}
				}
			}

		}
	}

	return game, nil
}
