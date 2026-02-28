package seeder

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	insertGame = `INSERT INTO game (description, created_at) VALUES (?, ?);`

	getPlayer    = `SELECT id FROM player WHERE name = ?;`
	insertPlayer = `INSERT INTO player (name) VALUES (?);`

	getDeck    = `SELECT id FROM deck WHERE player_id = ? AND commander = ?;`
	insertDeck = `INSERT INTO deck (player_id, commander) VALUES (?, ?);`

	insertGameResult = `INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES (?, ?, ?, ?);`

	getRoleByName     = `SELECT id FROM user_role WHERE name = ?;`
	getUserByPlayerID = `SELECT id FROM user WHERE player_id = ?;`
	insertUser        = `INSERT INTO user (player_id, role_id) VALUES (?, ?);`
)

type Seeder struct {
	client *lib.DBClient
	log    *zap.Logger
}

func NewSeeder(log *zap.Logger, client *lib.DBClient) *Seeder {
	return &Seeder{
		client: client,
		log:    log,
	}
}

func (s *Seeder) Run(ctx context.Context) error {
	s.log.Info("Running Data Seeder...")

	data, err := os.ReadFile("./data/gameInfos.json")
	if err != nil {
		return fmt.Errorf("failed to read game info json file: %w", err)
	}

	var games []Game
	if err := json.Unmarshal(data, &games); err != nil {
		return fmt.Errorf("failed to unmarshal game info: %w", err)
	}

	s.log.Info("Seeding Games", zap.Int("Count", len(games)))

	for i, game := range games {
		logger := s.log.With(zap.Any("Game", game))

		// Make a game record
		gameID, err := s.insertGame(ctx, i+1, game.Date)
		if err != nil {
			logger.Error("Error inserting game record", zap.Error(err))
			return err
		}

		for _, result := range game.Results {
			logger = logger.With(zap.Any("Result", result))

			// Get or insert the player
			playerID, err := s.getOrInsertPlayer(ctx, result.Player)
			if err != nil {
				logger.Error("Error getting or inserting player", zap.Error(err))
				return err
			}

			if err = s.getOrInsertUser(ctx, playerID); err != nil {
				logger.Error("Error getting or inserting user", zap.Error(err))
				return err
			}

			// Get or insert the deck
			deckID, err := s.getOrInsertDeck(ctx, playerID, result.Commander)
			if err != nil {
				logger.Error("Error getting or inserting deck", zap.Error(err))
				return err
			}

			// Create the game result
			if err := s.insertGameResult(ctx, gameID, deckID, result.Place, result.Kills); err != nil {
				logger.Error("Error inserting game result", zap.Error(err))
				return err
			}
		}
	}

	s.log.Info("Games seeded")

	return nil
}

func (s *Seeder) insertGame(ctx context.Context, count int, date time.Time) (int64, error) {
	result, err := s.client.Db.ExecContext(
		ctx,
		insertGame,
		fmt.Sprintf("Game %v", count),
		date,
	)
	if err != nil {
		return -1, fmt.Errorf("failed to insert game record: %w", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to get last inserted ID for new game record: %w", err)
	}

	return lastId, nil
}

func (s *Seeder) getOrInsertPlayer(ctx context.Context, name string) (int64, error) {
	var id int64
	if err := s.client.Db.QueryRowContext(ctx, getPlayer, name).Scan(&id); errors.Is(err, sql.ErrNoRows) {
		return s.insertPlayer(ctx, name)
	} else if err != nil {
		return -1, fmt.Errorf("failed to get player ID: %w", err)
	}

	return id, nil
}

func (s *Seeder) insertPlayer(ctx context.Context, name string) (int64, error) {
	result, err := s.client.Db.ExecContext(ctx, insertPlayer, name)
	if err != nil {
		return -1, fmt.Errorf("failed to insert player record: %w", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to get last inserted ID for new player record: %w", err)
	}

	return lastId, nil
}

func (s *Seeder) getOrInsertDeck(ctx context.Context, playerID int64, commander string) (int64, error) {
	var id int64
	if err := s.client.Db.QueryRowContext(ctx, getDeck, playerID, commander).Scan(&id); errors.Is(err, sql.ErrNoRows) {
		return s.insertDeck(ctx, playerID, commander)
	} else if err != nil {
		return -1, fmt.Errorf("failed to get deck ID: %w", err)
	}

	return id, nil
}

func (s *Seeder) insertDeck(ctx context.Context, playerID int64, commander string) (int64, error) {
	result, err := s.client.Db.ExecContext(ctx, insertDeck, playerID, commander)
	if err != nil {
		return -1, fmt.Errorf("failed to insert deck record: %w", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to get last inserted ID for new deck record: %w", err)
	}

	return lastId, nil
}

func (s *Seeder) insertGameResult(ctx context.Context, gameID, deckID int64, place, kills int) error {
	if _, err := s.client.Db.ExecContext(ctx, insertGameResult, gameID, deckID, place, kills); err != nil {
		return fmt.Errorf("failed to insert game result record: %w", err)
	}

	return nil
}

func (s *Seeder) getOrInsertUser(ctx context.Context, playerID int64) error {
	var roleID int64
	if err := s.client.Db.QueryRowContext(ctx, getRoleByName, "player").Scan(&roleID); err != nil {
		return fmt.Errorf("failed to get player role ID: %w", err)
	}

	var id int64
	if err := s.client.Db.QueryRowContext(ctx, getUserByPlayerID, playerID).Scan(&id); errors.Is(err, sql.ErrNoRows) {
		if _, err := s.client.Db.ExecContext(ctx, insertUser, playerID, roleID); err != nil {
			return fmt.Errorf("failed to insert user record: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to get user for player_id %d: %w", playerID, err)
	}

	return nil
}
