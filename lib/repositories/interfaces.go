package repositories

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib/repositories/commander"
	"github.com/m-sharp/edh-tracker/lib/repositories/deck"
	"github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	"github.com/m-sharp/edh-tracker/lib/repositories/format"
	"github.com/m-sharp/edh-tracker/lib/repositories/game"
	"github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	"github.com/m-sharp/edh-tracker/lib/repositories/player"
	"github.com/m-sharp/edh-tracker/lib/repositories/pod"
	"github.com/m-sharp/edh-tracker/lib/repositories/user"
)

type PlayerRepository interface {
	GetAll(ctx context.Context) ([]player.Model, error)
	GetById(ctx context.Context, playerID int) (*player.Model, error)
	GetByName(ctx context.Context, name string) (*player.Model, error)
	GetByNames(ctx context.Context, names []string) ([]player.Model, error)
	Add(ctx context.Context, name string) (int, error)
	BulkAdd(ctx context.Context, names []string) ([]player.Model, error)
	SoftDelete(ctx context.Context, id int) error
}

type DeckRepository interface {
	GetAll(ctx context.Context) ([]deck.Model, error)
	GetAllForPlayer(ctx context.Context, playerID int) ([]deck.Model, error)
	GetById(ctx context.Context, deckID int) (*deck.Model, error)
	Add(ctx context.Context, playerID int, name string, formatID int) (int, error)
	BulkAdd(ctx context.Context, decks []deck.Model) ([]deck.Model, error)
	Retire(ctx context.Context, deckID int) error
	SoftDelete(ctx context.Context, id int) error
}

type GameRepository interface {
	GetAllByPod(ctx context.Context, podID int) ([]game.Model, error)
	GetAllByDeck(ctx context.Context, deckID int) ([]game.Model, error)
	GetById(ctx context.Context, gameID int) (*game.Model, error)
	Add(ctx context.Context, description string, podID, formatID int) (int, error)
	BulkAdd(ctx context.Context, games []game.Model) ([]int, error)
	SoftDelete(ctx context.Context, id int) error
}

type GameResultRepository interface {
	GetByGameId(ctx context.Context, gameID int) ([]gameResult.Model, error)
	GetStatsForPlayer(ctx context.Context, playerID int) (*gameResult.Aggregate, error)
	GetStatsForDeck(ctx context.Context, deckID int) (*gameResult.Aggregate, error)
	BulkAdd(ctx context.Context, results []gameResult.Model) error
	SoftDelete(ctx context.Context, id int) error
}

type PodRepository interface {
	GetAll(ctx context.Context) ([]pod.Model, error)
	GetByID(ctx context.Context, podID int) (*pod.Model, error)
	GetByPlayerID(ctx context.Context, playerID int) ([]pod.Model, error)
	GetByName(ctx context.Context, name string) (*pod.Model, error)
	GetIDsByPlayerID(ctx context.Context, playerID int) ([]int, error)
	Add(ctx context.Context, name string) (int, error)
	BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error
	AddPlayerToPod(ctx context.Context, podID, playerID int) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*user.Model, error)
	GetByPlayerID(ctx context.Context, playerID int) (*user.Model, error)
	GetByOAuth(ctx context.Context, provider, subject string) (*user.Model, error)
	GetRoleByName(ctx context.Context, name string) (*user.RoleModel, error)
	Add(ctx context.Context, playerID, roleID int) (int, error)
	AddWithOAuth(ctx context.Context, playerID, roleID int, provider, subject, email, displayName, avatarURL string) (int, error)
	// CreatePlayerAndUser atomically inserts a player row and a linked user row in one transaction.
	// Returns the created user Model.
	CreatePlayerAndUser(ctx context.Context, playerName string, roleID int, provider, subject, email, displayName, avatarURL string) (*user.Model, error)
	BulkAdd(ctx context.Context, playerIDs []int, roleID int) error
	SoftDelete(ctx context.Context, id int) error
}

type FormatRepository interface {
	GetAll(ctx context.Context) ([]format.Model, error)
	GetById(ctx context.Context, id int) (*format.Model, error)
	GetByName(ctx context.Context, name string) (*format.Model, error)
}

type CommanderRepository interface {
	GetById(ctx context.Context, id int) (*commander.Model, error)
	GetByName(ctx context.Context, name string) (*commander.Model, error)
	GetByNames(ctx context.Context, names []string) ([]commander.Model, error)
	Add(ctx context.Context, name string) (int, error)
	BulkAdd(ctx context.Context, names []string) ([]commander.Model, error)
}

type DeckCommanderRepository interface {
	GetByDeckId(ctx context.Context, deckID int) (*deckCommander.Model, error)
	Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error)
	BulkAdd(ctx context.Context, entries []deckCommander.Model) error
}
