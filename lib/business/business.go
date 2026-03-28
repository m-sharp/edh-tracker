package business

import (
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
	"github.com/m-sharp/edh-tracker/lib/business/commander"
	"github.com/m-sharp/edh-tracker/lib/business/deck"
	"github.com/m-sharp/edh-tracker/lib/business/format"
	"github.com/m-sharp/edh-tracker/lib/business/game"
	"github.com/m-sharp/edh-tracker/lib/business/gameResult"
	"github.com/m-sharp/edh-tracker/lib/business/player"
	"github.com/m-sharp/edh-tracker/lib/business/pod"
	"github.com/m-sharp/edh-tracker/lib/business/user"
	"github.com/m-sharp/edh-tracker/lib/repositories"
)

type Business struct {
	Players     player.Functions
	Decks       deck.Functions
	Games       game.Functions
	GameResults gameResult.Functions
	Formats     format.Functions
	Commanders  commander.Functions
	Pods        pod.Functions
	Users       user.Functions
}

func NewBusiness(log *zap.Logger, r *repositories.Repositories, client *lib.DBClient) *Business {
	// Build leaf functions first (no cross-domain deps).
	getFormat := format.GetByID(r.Formats)
	getCommanderName := commander.GetCommanderName(r.Commanders)

	// Build game result functions.
	enrichGameResults := gameResult.EnrichModels()

	return &Business{
		Players: player.Functions{
			GetAll:        player.GetAll(r.Players, r.GameResults, r.Pods),
			GetAllByPod:   player.GetAllByPod(r.Players, r.GameResults, r.Pods, r.PlayerPodRoles),
			GetByID:       player.GetByID(r.Players, r.GameResults, r.Pods),
			Create:        player.Create(r.Players),
			Update:        player.Update(r.Players),
			GetPlayerName: player.GetPlayerName(r.Players),
		},
		Decks: deck.Functions{
			GetAllForPlayer:         deck.GetAllForPlayer(r.Decks, r.GameResults),
			GetAllByPod:             deck.GetAllByPod(r.Decks, r.Pods, r.GameResults),
			GetAllByPodPaginated:    deck.GetAllByPodPaginated(r.Decks, r.GameResults),
			GetAllByPlayerPaginated: deck.GetAllByPlayerPaginated(r.Decks, r.GameResults),
			GetByID:                 deck.GetByID(r.Decks, r.GameResults),
			Create:                  deck.Create(r.Decks, r.DeckCommanders, getFormat),
			Update:                  deck.Update(r.Decks, r.DeckCommanders),
			SoftDelete:              deck.SoftDelete(r.Decks),
			Retire:                  deck.Retire(r.Decks),
			GetDeckName:             deck.GetDeckName(r.Decks),
		},
		Games: game.Functions{
			GetAllByPod:               game.GetAllByPod(log, r.Games, enrichGameResults),
			GetAllByDeck:              game.GetAllByDeck(log, r.Games, enrichGameResults),
			GetAllByPlayer:            game.GetAllByPlayer(log, r.Games, enrichGameResults),
			GetAllByPodPaginated:      game.GetAllByPodPaginated(log, r.Games, enrichGameResults),
			GetAllByDeckPaginated:     game.GetAllByDeckPaginated(log, r.Games, enrichGameResults),
			GetAllByPlayerIDPaginated: game.GetAllByPlayerIDPaginated(log, r.Games, enrichGameResults),
			GetByID:                   game.GetByID(log, r.Games, enrichGameResults),
			Create:                    game.Create(r.Games, r.GameResults, r.Decks, getFormat, client),
			Update:                    game.Update(r.Games),
			SoftDelete:                game.SoftDelete(r.Games),
			AddResult:                 game.AddResult(r.GameResults),
			UpdateResult:              game.UpdateResult(r.GameResults),
			DeleteResult:              game.DeleteResult(r.GameResults),
		},
		GameResults: gameResult.Functions{
			GetByGameID:        gameResult.GetByGameID(r.GameResults),
			GetGameIDForResult: gameResult.GetGameIDForResult(r.GameResults),
			EnrichModels:       enrichGameResults,
		},
		Formats: format.Functions{
			GetAll:  format.GetAll(r.Formats),
			GetByID: getFormat,
		},
		Commanders: commander.Functions{
			GetAll:           commander.GetAll(r.Commanders),
			GetByID:          commander.GetByID(r.Commanders),
			Create:           commander.Create(r.Commanders),
			GetCommanderName: getCommanderName,
		},
		Pods: pod.Functions{
			GetByID:             pod.GetByID(r.Pods),
			GetByPlayerID:       pod.GetByPlayerID(r.Pods),
			Create:              pod.Create(r.Pods, r.PlayerPodRoles, client),
			AddPlayer:           pod.AddPlayer(r.Pods, r.PlayerPodRoles),
			GetRole:             pod.GetRole(r.PlayerPodRoles),
			PromoteToManager:    pod.PromoteToManager(r.PlayerPodRoles),
			GenerateInvite:      pod.GenerateInvite(r.PodInvites),
			JoinByInvite:        pod.JoinByInvite(r.PodInvites, r.Pods, r.PlayerPodRoles),
			Leave:               pod.Leave(r.Pods, r.PlayerPodRoles),
			SoftDelete:          pod.SoftDelete(r.Pods),
			Update:              pod.Update(r.Pods),
			GetMembersWithRoles: pod.GetMembersWithRoles(r.PlayerPodRoles),
			RemovePlayer:        pod.RemovePlayer(r.Pods, r.PlayerPodRoles),
		},
		Users: user.Functions{
			GetByID:         user.GetByID(r.Users),
			GetByPlayerID:   user.GetByPlayerID(r.Users),
			Create:          user.Create(r.Users),
			GetByOAuth:      user.GetByOAuth(r.Users),
			CreateWithOAuth: user.CreateWithOAuth(r.Players, r.Users, client),
			GetByEmail:      user.GetByEmail(r.Users),
			LinkOAuth:       user.LinkOAuth(r.Users),
		},
	}
}
