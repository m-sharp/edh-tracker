package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib/repositories"
	"github.com/m-sharp/edh-tracker/lib/repositories/deck"
	"github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	"github.com/m-sharp/edh-tracker/lib/repositories/game"
	"github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	"github.com/m-sharp/edh-tracker/lib/repositories/user"
)

const DefaultPodName = "OG EDH Pod"

type Seeder struct {
	log   *zap.Logger
	repos *repositories.Repositories
}

func NewSeeder(log *zap.Logger, repos *repositories.Repositories) *Seeder {
	return &Seeder{
		log:   log.Named("Seeder"),
		repos: repos,
	}
}

// deckEntry holds the unique player+commander+format combinations found across all games.
type deckEntry struct {
	playerName    string
	commanderName string
	formatID      int
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

	// Pre-load format IDs from the database
	formats, err := s.repos.Formats.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load formats: %w", err)
	}
	formatIDs := make(map[string]int)
	for _, f := range formats {
		formatIDs[f.Name] = f.ID
	}

	// Load game data from JSON
	data, err := os.ReadFile("./data/gameInfos.json")
	if err != nil {
		return fmt.Errorf("failed to read game info json file: %w", err)
	}
	var games []Game
	if err = json.Unmarshal(data, &games); err != nil {
		return fmt.Errorf("failed to unmarshal game info: %w", err)
	}
	s.log.Info("Seeding Games", zap.Int("Count", len(games)))

	// Look up the player role once
	role, err := s.repos.Users.GetRoleByName(ctx, user.RolePlayer)
	if err != nil {
		return fmt.Errorf("failed to get player role: %w", err)
	}

	// Create the default pod
	podID, err := s.repos.Pods.Add(ctx, DefaultPodName)
	if err != nil {
		return fmt.Errorf("failed to create default pod: %w", err)
	}

	// Pre-processing pass: collect unique players, commanders, and deck entries from game data
	playerNames, commanderNames, deckEntries, err := s.collectEntities(games, formatIDs)
	if err != nil {
		return err
	}

	// Seed players, create user accounts, and add them to the pod
	playerIDs, err := s.seedPlayersAndUsers(ctx, playerNames, podID, role.ID)
	if err != nil {
		return err
	}

	// Seed commanders
	commanderIDs, err := s.seedCommanders(ctx, commanderNames)
	if err != nil {
		return err
	}

	// Seed decks and deck_commander associations
	deckIDs, err := s.seedDecks(ctx, deckEntries, playerIDs, commanderIDs)
	if err != nil {
		return err
	}

	// Seed games and their results
	if err = s.seedGames(ctx, games, formatIDs, playerIDs, deckIDs, podID); err != nil {
		return err
	}

	s.log.Info("Seeding complete", zap.Int("Games", len(games)))
	return nil
}

// collectEntities performs a pre-processing pass over the raw game data to collect
// unique player names, commander names, and deck entries (player+commander+format combos).
func (s *Seeder) collectEntities(games []Game, formatIDs map[string]int) (playerNames, commanderNames []string, entries []deckEntry, err error) {
	playerNameSet := map[string]struct{}{}
	commanderNameSet := map[string]struct{}{}
	deckKeySet := map[string]struct{}{}

	for i, gameToSeed := range games {
		formatID, ok := formatIDs[gameToSeed.Format]
		if !ok {
			return nil, nil, nil, fmt.Errorf("unknown format %q in gameToSeed %d", gameToSeed.Format, i+1)
		}
		for _, result := range gameToSeed.Results {
			playerNameSet[result.Player] = struct{}{}
			commanderNameSet[result.Name] = struct{}{}
			key := result.Player + ":" + result.Name
			if _, exists := deckKeySet[key]; !exists {
				deckKeySet[key] = struct{}{}
				entries = append(entries, deckEntry{
					playerName:    result.Player,
					commanderName: result.Name,
					formatID:      formatID,
				})
			}
		}
	}

	playerNames = make([]string, 0, len(playerNameSet))
	for name := range playerNameSet {
		playerNames = append(playerNames, name)
	}
	commanderNames = make([]string, 0, len(commanderNameSet))
	for name := range commanderNameSet {
		commanderNames = append(commanderNames, name)
	}

	return playerNames, commanderNames, entries, nil
}

// seedPlayersAndUsers bulk-inserts players, creates their user accounts, and adds them to the pod.
// Returns a name→ID map for use in subsequent seeding steps.
func (s *Seeder) seedPlayersAndUsers(ctx context.Context, playerNames []string, podID, roleID int) (map[string]int, error) {
	players, err := s.repos.Players.BulkAdd(ctx, playerNames)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk add players: %w", err)
	}

	playerIDs := make(map[string]int, len(players))
	playerIDSlice := make([]int, len(players))
	for i, p := range players {
		playerIDs[p.Name] = p.ID
		playerIDSlice[i] = p.ID
	}

	if err = s.repos.Users.BulkAdd(ctx, playerIDSlice, roleID); err != nil {
		return nil, fmt.Errorf("failed to bulk add users: %w", err)
	}

	if err = s.repos.Pods.BulkAddPlayers(ctx, podID, playerIDSlice); err != nil {
		return nil, fmt.Errorf("failed to bulk add players to pod: %w", err)
	}

	if err = s.repos.PlayerPodRoles.BulkAdd(ctx, podID, playerIDSlice, "member"); err != nil {
		return nil, fmt.Errorf("failed to bulk add player pod roles: %w", err)
	}

	mikeID, ok := playerIDs["Mike"]
	if ok {
		if err = s.repos.PlayerPodRoles.SetRole(ctx, podID, mikeID, "manager"); err != nil {
			return nil, fmt.Errorf("failed to set Mike as pod manager: %w", err)
		}
	}

	return playerIDs, nil
}

// seedCommanders bulk-inserts commanders and returns a name→ID map.
func (s *Seeder) seedCommanders(ctx context.Context, commanderNames []string) (map[string]int, error) {
	commanders, err := s.repos.Commanders.BulkAdd(ctx, commanderNames)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk add commanders: %w", err)
	}

	commanderIDs := make(map[string]int, len(commanders))
	for _, c := range commanders {
		commanderIDs[c.Name] = c.ID
	}

	return commanderIDs, nil
}

// seedDecks bulk-inserts decks and their commander associations.
// Returns a "playerID:deckName"→deckID map for use when building game results.
func (s *Seeder) seedDecks(ctx context.Context, entries []deckEntry, playerIDs, commanderIDs map[string]int) (map[string]int, error) {
	// Build the []deck.Model slice for bulk insertion
	deckSlice := make([]deck.Model, len(entries))
	for i, de := range entries {
		deckSlice[i] = deck.Model{
			PlayerID: playerIDs[de.playerName],
			Name:     de.commanderName,
			FormatID: de.formatID,
		}
	}

	insertedDecks, err := s.repos.Decks.BulkAdd(ctx, deckSlice)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk add decks: %w", err)
	}

	deckIDs := make(map[string]int, len(insertedDecks))
	for _, d := range insertedDecks {
		key := fmt.Sprintf("%d:%s", d.PlayerID, d.Name)
		deckIDs[key] = d.ID
	}

	// Build the []deckCommander.Model slice and insert associations
	dcSlice := make([]deckCommander.Model, len(entries))
	for i, de := range entries {
		playerID := playerIDs[de.playerName]
		deckKey := fmt.Sprintf("%d:%s", playerID, de.commanderName)
		dcSlice[i] = deckCommander.Model{
			DeckID:      deckIDs[deckKey],
			CommanderID: commanderIDs[de.commanderName],
		}
	}

	if err = s.repos.DeckCommanders.BulkAdd(ctx, dcSlice); err != nil {
		return nil, fmt.Errorf("failed to bulk add deck_commanders: %w", err)
	}

	return deckIDs, nil
}

// seedGames bulk-inserts game records and then their results separately.
func (s *Seeder) seedGames(ctx context.Context, games []Game, formatIDs, playerIDs, deckIDs map[string]int, podID int) error {
	// Phase A: build and insert game records
	gameModels := make([]game.Model, len(games))
	for i, g := range games {
		gameModels[i] = game.Model{
			Description: fmt.Sprintf("Game %d", i+1),
			PodID:       podID,
			FormatID:    formatIDs[g.Format],
		}
	}
	gameIDs, err := s.repos.Games.BulkAdd(ctx, gameModels)
	if err != nil {
		return fmt.Errorf("failed to bulk add games: %w", err)
	}

	// Phase B: build and insert game result records
	var allResults []gameResult.Model
	for i, g := range games {
		for _, result := range g.Results {
			playerID := playerIDs[result.Player]
			deckKey := fmt.Sprintf("%d:%s", playerID, result.Name)
			allResults = append(allResults, gameResult.Model{
				GameID:    gameIDs[i],
				DeckID:    deckIDs[deckKey],
				Place:     result.Place,
				KillCount: result.Kills,
			})
		}
	}
	if err = s.repos.GameResults.BulkAdd(ctx, allResults); err != nil {
		return fmt.Errorf("failed to bulk add game results: %w", err)
	}

	return nil
}
