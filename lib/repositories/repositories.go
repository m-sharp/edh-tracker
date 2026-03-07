package repositories

import (
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
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

// Compile-time interface satisfaction checks.
var (
	_ PlayerRepository        = (*player.Repository)(nil)
	_ DeckRepository          = (*deck.Repository)(nil)
	_ GameRepository          = (*game.Repository)(nil)
	_ GameResultRepository    = (*gameResult.Repository)(nil)
	_ PodRepository           = (*pod.Repository)(nil)
	_ UserRepository          = (*user.Repository)(nil)
	_ FormatRepository        = (*format.Repository)(nil)
	_ CommanderRepository     = (*commander.Repository)(nil)
	_ DeckCommanderRepository = (*deckCommander.Repository)(nil)
	_ PlayerPodRoleRepository = (*playerPodRole.Repository)(nil)
	_ PodInviteRepository     = (*podInvite.Repository)(nil)
)

type Repositories struct {
	Players        *player.Repository
	Decks          *deck.Repository
	Games          *game.Repository
	GameResults    *gameResult.Repository
	Pods           *pod.Repository
	Users          *user.Repository
	Formats        *format.Repository
	Commanders     *commander.Repository
	DeckCommanders *deckCommander.Repository
	PlayerPodRoles *playerPodRole.Repository
	PodInvites     *podInvite.Repository
}

func New(_ *zap.Logger, client *lib.DBClient) *Repositories {
	return &Repositories{
		Players:        player.NewRepository(client),
		Decks:          deck.NewRepository(client),
		Games:          game.NewRepository(client),
		GameResults:    gameResult.NewRepository(client),
		Pods:           pod.NewRepository(client),
		Users:          user.NewRepository(client),
		Formats:        format.NewRepository(client),
		Commanders:     commander.NewRepository(client),
		DeckCommanders: deckCommander.NewRepository(client),
		PlayerPodRoles: playerPodRole.NewRepository(client),
		PodInvites:     podInvite.NewRepository(client),
	}
}
