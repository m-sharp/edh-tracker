package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/models"
)

const DefaultPodName = "OG EDH Pod"

type Seeder struct {
	log       *zap.Logger
	repos     *models.Repositories
	playerIDs map[string]int
	deckIDs   map[string]int
}

func NewSeeder(log *zap.Logger, repos *models.Repositories) *Seeder {
	return &Seeder{
		log:       log.Named("Seeder"),
		repos:     repos,
		playerIDs: map[string]int{},
		deckIDs:   map[string]int{},
	}
}

func (s *Seeder) Run(ctx context.Context) error {
	s.log.Info("Running Data Seeder...")

	// Guard against re-runs: if the default pod already exists, seed data is already present
	existing, err := s.repos.Pods.GetByName(ctx, DefaultPodName)
	if err != nil {
		return fmt.Errorf("failed to check for existing seed data: %w", err)
	}
	if existing != nil {
		s.log.Warn("Seed data already exists, skipping seeder", zap.String("Pod", DefaultPodName))
		return nil
	}

	data, err := os.ReadFile("./data/gameInfos.json")
	if err != nil {
		return fmt.Errorf("failed to read game info json file: %w", err)
	}

	var games []Game
	if err = json.Unmarshal(data, &games); err != nil {
		return fmt.Errorf("failed to unmarshal game info: %w", err)
	}

	s.log.Info("Seeding Games", zap.Int("Count", len(games)))

	// Look up the player role once and cache the ID
	role, err := s.repos.Users.GetRoleByName(ctx, models.RolePlayer)
	if err != nil {
		return fmt.Errorf("failed to get player role: %w", err)
	}
	roleID := role.ID

	// Create the default pod
	podID, err := s.repos.Pods.Add(ctx, DefaultPodName)
	if err != nil {
		return fmt.Errorf("failed to create default pod: %w", err)
	}

	for i, game := range games {
		var results []models.GameResult

		for _, result := range game.Results {
			playerID, err := s.getOrCreatePlayer(ctx, result.Player, podID, roleID)
			if err != nil {
				return fmt.Errorf("failed to get or create player %q: %w", result.Player, err)
			}

			deckID, err := s.getOrCreateDeck(ctx, playerID, result.Commander)
			if err != nil {
				return fmt.Errorf("failed to get or create deck %q for player %d: %w", result.Commander, playerID, err)
			}

			results = append(results, models.GameResult{
				DeckId: deckID,
				Place:  result.Place,
				Kills:  result.Kills,
			})
		}

		description := fmt.Sprintf("Game %d", i+1)
		if err = s.repos.Games.Add(ctx, description, podID, results...); err != nil {
			return fmt.Errorf("failed to insert game %d: %w", i+1, err)
		}
	}

	s.log.Info("Seeding complete", zap.Int("Games", len(games)))
	return nil
}

func (s *Seeder) getOrCreatePlayer(ctx context.Context, name string, podID, roleID int) (int, error) {
	if id, ok := s.playerIDs[name]; ok {
		return id, nil
	}

	playerID, err := s.repos.Players.Add(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("failed to add player %q: %w", name, err)
	}

	if _, err = s.repos.Users.Add(ctx, playerID, roleID); err != nil {
		return 0, fmt.Errorf("failed to add user for player %d: %w", playerID, err)
	}

	if err = s.repos.Pods.AddPlayerToPod(ctx, podID, playerID); err != nil {
		return 0, fmt.Errorf("failed to add player %d to pod %d: %w", playerID, podID, err)
	}

	s.playerIDs[name] = playerID
	return playerID, nil
}

func (s *Seeder) getOrCreateDeck(ctx context.Context, playerID int, commander string) (int, error) {
	key := fmt.Sprintf("%d:%s", playerID, commander)
	if id, ok := s.deckIDs[key]; ok {
		return id, nil
	}

	deckID, err := s.repos.Decks.Add(ctx, playerID, commander)
	if err != nil {
		return 0, fmt.Errorf("failed to add deck %q for player %d: %w", commander, playerID, err)
	}

	s.deckIDs[key] = deckID
	return deckID, nil
}
