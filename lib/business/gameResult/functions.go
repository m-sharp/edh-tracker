package gameResult

import (
	"context"
	"fmt"

	repos "github.com/m-sharp/edh-tracker/lib/repositories"
	gameResultRepo "github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	"github.com/m-sharp/edh-tracker/lib/utils"
)

func GetByGameID(gameResultRepo repos.GameResultRepository) GetByGameIDFunc {
	enrich := EnrichModels()
	return func(ctx context.Context, gameID int) ([]Entity, error) {
		resultModels, err := gameResultRepo.GetByGameID(ctx, gameID)
		if err != nil {
			return nil, fmt.Errorf("failed to get results for game %d: %w", gameID, err)
		}
		return enrich(ctx, resultModels)
	}
}

// TODO: Not a huge fan of this, feeds back into the question of do we need a view for points
func EnrichModels() EnrichModelsFunc {
	return func(ctx context.Context, models []gameResultRepo.Model) ([]Entity, error) {
		numPlayers := len(models)
		results := make([]Entity, 0, len(models))
		for _, r := range models {
			entity := Entity{
				ID:       r.ID,
				GameID:   r.GameID,
				DeckID:   r.DeckID,
				PlayerID: r.Deck.PlayerID,
				DeckName: r.Deck.Name,
				Place:    r.Place,
				Kills:    r.KillCount,
				Points:   utils.GetPointsForPlace(r.KillCount, r.Place, numPlayers),
			}
			if r.Deck.Commander != nil {
				name := r.Deck.Commander.Commander.Name
				entity.CommanderName = &name
				if r.Deck.Commander.PartnerCommander != nil {
					partnerName := r.Deck.Commander.PartnerCommander.Name
					entity.PartnerCommanderName = &partnerName
				}
			}
			results = append(results, entity)
		}
		return results, nil
	}
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
