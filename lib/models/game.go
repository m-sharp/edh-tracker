package models

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

const (
	GetAllGames      = `SELECT id, description, ctime FROM game;`
	GetGameByID      = `SELECT id, description, ctime FROM game WHERE id = ?;`
	GetGamesByDeckId = `SELECT game.id, game.description, game.ctime
								FROM (game INNER JOIN game_result on game.id = game_result.game_id)
							  WHERE game_result.deck_id = ?;`
	GetGameResultsByGameID = `SELECT game_result.id, game_result.game_id, game_result.deck_id, game_result.place, game_result.kill_count, deck.commander
								FROM (game_result INNER JOIN deck on game_result.deck_id = deck.id)
							  WHERE game_id = ?;`

	InsertGame       = `INSERT INTO game (description) VALUES (?);`
	InsertGameResult = `INSERT INTO game_result (game_id, deck_id, place, kill_count) VALUES (?, ?, ?, ?);`

	gameResultValidationErr = "invalid Game Result: %s"
)

type Game struct {
	Id          int       `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"ctime" db:"ctime"`
}

type GameResult struct {
	Id        int    `json:"id" db:"id"`
	GameId    int    `json:"game_id" db:"game_id"`
	DeckId    int    `json:"deck_id" db:"deck_id"`
	Commander string `json:"commander" db:"commander"`
	Place     int    `json:"place" db:"place"`
	Kills     int    `json:"kill_count" db:"kill_count"`
	Points    int    `json:"points"`
}

func (g *GameResult) Validate() error {
	if g.GameId == 0 {
		return fmt.Errorf(gameResultValidationErr, "missing GameId")
	}
	if g.DeckId == 0 {
		return fmt.Errorf(gameResultValidationErr, "missing DeckId")
	}
	if g.Place == 0 {
		return fmt.Errorf(gameResultValidationErr, "missing Place")
	}
	if g.Place < 1 {
		return fmt.Errorf(gameResultValidationErr, "Place cannot be less than 1")
	}
	if g.Kills < 0 {
		return fmt.Errorf(gameResultValidationErr, "Kills cannot be less than 0")
	}

	return nil
}

type GameDetails struct {
	Game
	Results []GameResult `json:"results"`
}

type GameProvider struct {
	log    *zap.Logger
	client *lib.DBClient
}

func NewGameProvider(log *zap.Logger, client *lib.DBClient) *GameProvider {
	return &GameProvider{
		log:    log.Named("GameProvider"),
		client: client,
	}
}

func (g *GameProvider) GetAll(ctx context.Context) ([]GameDetails, error) {
	var games []Game
	if err := g.client.Db.SelectContext(ctx, &games, GetAllGames); err != nil {
		return nil, fmt.Errorf("failed to get Game records: %w", err)
	}

	if games == nil {
		return []GameDetails{}, nil
	}

	var details []GameDetails
	for _, game := range games {
		results, err := g.getGameResults(ctx, game.Id)
		if err != nil {
			g.log.Warn("Failed to get game results for game, dropping from results", zap.Any("Game", game))
			continue
		}

		details = append(details, GameDetails{Game: game, Results: results})
	}

	return details, nil
}

func (g *GameProvider) GetAllByDeck(ctx context.Context, deckId int) ([]GameDetails, error) {
	var games []Game
	if err := g.client.Db.SelectContext(ctx, &games, GetGamesByDeckId, deckId); err != nil {
		return nil, fmt.Errorf("failed to get Game records: %w", err)
	}

	if games == nil {
		return []GameDetails{}, nil
	}

	var details []GameDetails
	for _, game := range games {
		results, err := g.getGameResults(ctx, game.Id)
		if err != nil {
			g.log.Warn("Failed to get game results for game, dropping from results", zap.Any("Game", game))
			continue
		}

		details = append(details, GameDetails{Game: game, Results: results})
	}

	return details, nil
}

func (g *GameProvider) GetGameById(ctx context.Context, gameId int) (*GameDetails, error) {
	var games []Game
	if err := g.client.Db.SelectContext(ctx, &games, GetGameByID, gameId); err != nil {
		return nil, fmt.Errorf("failed to get Game record for id %d: %w", gameId, err)
	}

	if len(games) == 0 || len(games) > 1 {
		return nil, fmt.Errorf(
			"unexpected number of games returned for ID %d: got %d, expected 1",
			gameId, len(games),
		)
	}

	game := games[0]
	results, err := g.getGameResults(ctx, game.Id)
	if err != nil {
		return nil, err
	}

	return &GameDetails{
		Game:    game,
		Results: results,
	}, nil
}

func (g *GameProvider) getGameResults(ctx context.Context, gameId int) ([]GameResult, error) {
	var results []GameResult
	if err := g.client.Db.SelectContext(ctx, &results, GetGameResultsByGameID, gameId); err != nil {
		return nil, fmt.Errorf("failed to get Game Results for Game %d: %w", gameId, err)
	}

	if results == nil {
		return nil, fmt.Errorf("failed to get Game Results for Game %d: no results found", gameId)
	}

	for i := 0; i < len(results); i++ {
		target := results[i]
		results[i].Points = getPointsForPlace(target.Kills, target.Place)
	}

	return results, nil
}

func (g *GameProvider) Add(ctx context.Context, description string, results ...GameResult) error {
	r, err := g.client.Db.ExecContext(ctx, InsertGame, description)
	if err != nil {
		return fmt.Errorf("failed to insert Game record: %w", err)
	}

	numAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
	}
	if numAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected by Game insert: got %d, expected 1", numAffected)
	}

	newId, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	for _, result := range results {
		r, err := g.client.Db.ExecContext(ctx, InsertGameResult, newId, result.DeckId, result.Place, result.Kills)
		if err != nil {
			return fmt.Errorf("failed to insert Game Result record: %w", err)
		}

		numAffected, err := r.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get number of rows affected by insert: %w", err)
		}
		if numAffected != 1 {
			return fmt.Errorf("unexpected number of rows affected by Game Result insert: got %d, expected 1", numAffected)
		}
	}

	return nil
}
