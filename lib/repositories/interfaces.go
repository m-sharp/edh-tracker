package repositories

import (
	"context"
	"time"

	"github.com/m-sharp/edh-tracker/lib/repositories/commander"
	"github.com/m-sharp/edh-tracker/lib/repositories/deck"
	"github.com/m-sharp/edh-tracker/lib/repositories/deckCommander"
	"github.com/m-sharp/edh-tracker/lib/repositories/format"
	"github.com/m-sharp/edh-tracker/lib/repositories/game"
	"github.com/m-sharp/edh-tracker/lib/repositories/gameResult"
	"github.com/m-sharp/edh-tracker/lib/repositories/player"
	"github.com/m-sharp/edh-tracker/lib/repositories/playerPodRole"
	"github.com/m-sharp/edh-tracker/lib/repositories/pod"
	"github.com/m-sharp/edh-tracker/lib/repositories/podInvite"
	"github.com/m-sharp/edh-tracker/lib/repositories/user"
)

type PlayerRepository interface {
	GetAll(ctx context.Context) ([]player.Model, error)
	GetById(ctx context.Context, playerID int) (*player.Model, error)
	GetByName(ctx context.Context, name string) (*player.Model, error)
	GetByNames(ctx context.Context, names []string) ([]player.Model, error)
	Add(ctx context.Context, name string) (int, error)
	BulkAdd(ctx context.Context, names []string) ([]player.Model, error)
	Update(ctx context.Context, playerID int, name string) error
	SoftDelete(ctx context.Context, id int) error
}

type DeckRepository interface {
	GetAll(ctx context.Context) ([]deck.Model, error)
	GetAllHydrated(ctx context.Context) ([]deck.Model, error)
	GetAllForPlayer(ctx context.Context, playerID int) ([]deck.Model, error)
	GetAllForPlayerHydrated(ctx context.Context, playerID int) ([]deck.Model, error)
	GetAllByPlayerIDs(ctx context.Context, playerIDs []int) ([]deck.Model, error)
	GetAllByPlayerIDsHydrated(ctx context.Context, playerIDs []int) ([]deck.Model, error)
	GetById(ctx context.Context, deckID int) (*deck.Model, error)
	GetByIDHydrated(ctx context.Context, deckID int) (*deck.Model, error)
	Add(ctx context.Context, playerID int, name string, formatID int) (int, error)
	BulkAdd(ctx context.Context, decks []deck.Model) ([]deck.Model, error)
	Update(ctx context.Context, deckID int, fields deck.UpdateFields) error
	Retire(ctx context.Context, deckID int) error
	SoftDelete(ctx context.Context, id int) error
}

type GameRepository interface {
	GetAllByPod(ctx context.Context, podID int) ([]game.Model, error)
	GetAllByPodWithResults(ctx context.Context, podID int) ([]game.Model, error)
	GetAllByDeck(ctx context.Context, deckID int) ([]game.Model, error)
	GetAllByDeckWithResults(ctx context.Context, deckID int) ([]game.Model, error)
	GetAllByPlayerID(ctx context.Context, playerID int) ([]game.Model, error)
	GetAllByPlayerWithResults(ctx context.Context, playerID int) ([]game.Model, error)
	GetById(ctx context.Context, gameID int) (*game.Model, error)
	GetByIDWithResults(ctx context.Context, gameID int) (*game.Model, error)
	Add(ctx context.Context, description string, podID, formatID int) (int, error)
	BulkAdd(ctx context.Context, games []game.Model) ([]int, error)
	Update(ctx context.Context, gameID int, description string) error
	SoftDelete(ctx context.Context, id int) error
}

type GameResultRepository interface {
	// TODO: If we're super efficient regardless of scale, why not always fetch the hydrated versions and drop the non-hydrated versions?
	GetByGameId(ctx context.Context, gameID int) ([]gameResult.Model, error)
	GetByGameIDWithDeckInfo(ctx context.Context, gameID int) ([]gameResult.Model, error)
	GetByID(ctx context.Context, resultID int) (*gameResult.Model, error)
	GetStatsForPlayer(ctx context.Context, playerID int) (*gameResult.Aggregate, error)
	GetStatsForDeck(ctx context.Context, deckID int) (*gameResult.Aggregate, error)
	Add(ctx context.Context, m gameResult.Model) (int, error)
	BulkAdd(ctx context.Context, results []gameResult.Model) error
	Update(ctx context.Context, resultID, place, killCount, deckID int) error
	SoftDelete(ctx context.Context, id int) error
}

type PodRepository interface {
	GetAll(ctx context.Context) ([]pod.Model, error)
	GetByID(ctx context.Context, podID int) (*pod.Model, error)
	GetByPlayerID(ctx context.Context, playerID int) ([]pod.Model, error)
	GetByName(ctx context.Context, name string) (*pod.Model, error)
	GetIDsByPlayerID(ctx context.Context, playerID int) ([]int, error)
	GetPlayerIDs(ctx context.Context, podID int) ([]int, error)
	Add(ctx context.Context, name string) (int, error)
	BulkAddPlayers(ctx context.Context, podID int, playerIDs []int) error
	AddPlayerToPod(ctx context.Context, podID, playerID int) error
	SoftDelete(ctx context.Context, podID int) error
	Update(ctx context.Context, podID int, name string) error
	RemovePlayer(ctx context.Context, podID, playerID int) error
}

type PlayerPodRoleRepository interface {
	GetRole(ctx context.Context, podID, playerID int) (*playerPodRole.Model, error)
	SetRole(ctx context.Context, podID, playerID int, role string) error
	GetMembersWithRoles(ctx context.Context, podID int) ([]playerPodRole.Model, error)
	BulkAdd(ctx context.Context, podID int, playerIDs []int, role string) error
}

type PodInviteRepository interface {
	GetByCode(ctx context.Context, code string) (*podInvite.Model, error)
	Add(ctx context.Context, podID, createdByPlayerID int, code string, expiresAt *time.Time) error
	IncrementUsedCount(ctx context.Context, code string) error
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
	DeleteByDeckID(ctx context.Context, deckID int) error
}
