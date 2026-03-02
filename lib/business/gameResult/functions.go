package gameResult

import (
	"context"
	"fmt"

	"github.com/m-sharp/edh-tracker/lib/business/deck"
	repos "github.com/m-sharp/edh-tracker/lib/repositories"
)

func GetByGameID(
	gameResultRepo repos.GameResultRepository,
	getDeckName deck.GetDeckNameFunc,
	getCommanderEntry deck.GetCommanderEntryFunc,
) GetByGameIDFunc {
	return func(ctx context.Context, gameID int) ([]Entity, error) {
		resultModels, err := gameResultRepo.GetByGameId(ctx, gameID)
		if err != nil {
			return nil, fmt.Errorf("failed to get results for game %d: %w", gameID, err)
		}

		deckNameCache := map[int]string{}

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

			entity := Entity{
				ID:       r.ID,
				GameID:   r.GameID,
				DeckID:   r.DeckID,
				DeckName: deckName,
				Place:    r.Place,
				Kills:    r.KillCount,
				Points:   GetPointsForPlace(r.KillCount, r.Place),
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
