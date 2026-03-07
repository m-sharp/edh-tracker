package gameResult

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib/business/deck"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	"github.com/m-sharp/edh-tracker/lib/utils"
)

func GetByGameID(
	gameResultRepo repos.GameResultRepository,
	getDeckName deck.GetDeckNameFunc,
	getCommanderEntry deck.GetCommanderEntryFunc,
	getPlayerIDForDeck deck.GetPlayerIDForDeckFunc,
) GetByGameIDFunc {
	return func(ctx context.Context, gameID int) ([]Entity, error) {
		resultModels, err := gameResultRepo.GetByGameId(ctx, gameID)
		if err != nil {
			return nil, fmt.Errorf("failed to get results for game %d: %w", gameID, err)
		}

		numPlayers := len(resultModels)
		deckNameCache := map[int]string{}
		playerIDCache := map[int]int{}

		results := make([]Entity, 0, len(resultModels))
		for _, r := range resultModels {
			var deckName string
			if name, ok := deckNameCache[r.DeckID]; ok {
				deckName = name
			} else {
				name, err := getDeckName(ctx, r.DeckID)
				if err != nil {
					return nil, err
				}
				deckNameCache[r.DeckID] = name
				deckName = name
			}

			playerID, err := cachedPlayerIDForDeck(ctx, r.DeckID, playerIDCache, getPlayerIDForDeck)
			if err != nil {
				return nil, err
			}

			entity := Entity{
				ID:       r.ID,
				GameID:   r.GameID,
				DeckID:   r.DeckID,
				PlayerID: playerID,
				DeckName: deckName,
				Place:    r.Place,
				Kills:    r.KillCount,
				Points:   utils.GetPointsForPlace(r.KillCount, r.Place, numPlayers),
			}

			commanders, err := getCommanderEntry(ctx, r.DeckID)
			if err != nil {
				return nil, fmt.Errorf("failed to get commander for deck %d: %w", r.DeckID, err)
			}
			if commanders != nil {
				name := commanders.CommanderName
				entity.CommanderName = &name
				entity.PartnerCommanderName = commanders.PartnerCommanderName
			}

			results = append(results, entity)
		}

		return results, nil
	}
}

func cachedPlayerIDForDeck(ctx context.Context, deckID int, cache map[int]int, getPlayerIDForDeck deck.GetPlayerIDForDeckFunc) (int, error) {
	if pid, ok := cache[deckID]; ok {
		return pid, nil
	}
	pid, err := getPlayerIDForDeck(ctx, deckID)
	if err != nil {
		return 0, fmt.Errorf("failed to get player for deck %d: %w", deckID, err)
	}
	cache[deckID] = pid
	return pid, nil
}

func GetGameIDForResult(gameResultRepo repos.GameResultRepository) GetGameIDForResultFunc {
	return func(ctx context.Context, resultID int) (int, error) {
		m, err := gameResultRepo.GetByID(ctx, resultID)
		if err != nil {
			return 0, fmt.Errorf("failed to look up result %d: %w", resultID, err)
		}
		if m == nil {
			return 0, fmt.Errorf("game result %d not found", resultID)
		}
		return m.GameID, nil
	}
}
