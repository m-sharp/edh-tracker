package business

import (
	"go.uber.org/zap"

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

func NewBusiness(log *zap.Logger, r *repositories.Repositories) *Business {
	// Build leaf functions first (no cross-domain deps).
	getPlayerName := player.GetPlayerName(r.Players)
	getFormat := format.GetByID(r.Formats)
	getCommanderName := commander.GetCommanderName(r.Commanders)

	// Build deck cross-domain functions.
	getCommanderEntry := deck.GetCommanderEntry(r.DeckCommanders, getCommanderName)
	getDeckName := deck.GetDeckName(r.Decks)

	// Build game result function.
	getGameResults := gameResult.GetByGameID(r.GameResults, getDeckName, getCommanderEntry)

	return &Business{
		Players: player.Functions{
			GetAll:        player.GetAll(r.Players, r.GameResults, r.Pods),
			GetByID:       player.GetByID(r.Players, r.GameResults, r.Pods),
			Create:        player.Create(r.Players),
			GetPlayerName: getPlayerName,
		},
		Decks: deck.Functions{
			GetAll:            deck.GetAll(r.Decks, r.GameResults, getPlayerName, getFormat, getCommanderEntry),
			GetAllForPlayer:   deck.GetAllForPlayer(r.Decks, r.GameResults, getPlayerName, getFormat, getCommanderEntry),
			GetByID:           deck.GetByID(r.Decks, r.GameResults, getPlayerName, getFormat, getCommanderEntry),
			Create:            deck.Create(r.Decks, r.DeckCommanders, getFormat),
			Retire:            deck.Retire(r.Decks),
			GetDeckName:       getDeckName,
			GetCommanderEntry: getCommanderEntry,
		},
		Games: game.Functions{
			GetAllByPod:  game.GetAllByPod(log, r.Games, getGameResults),
			GetAllByDeck: game.GetAllByDeck(log, r.Games, getGameResults),
			GetByID:      game.GetByID(log, r.Games, getGameResults),
			Create:       game.Create(log, r.Games, r.GameResults, r.Decks, getFormat),
		},
		GameResults: gameResult.Functions{
			GetByGameID: getGameResults,
		},
		Formats: format.Functions{
			GetAll:  format.GetAll(r.Formats),
			GetByID: getFormat,
		},
		Commanders: commander.Functions{
			GetByID:          commander.GetByID(r.Commanders),
			Create:           commander.Create(r.Commanders),
			GetCommanderName: getCommanderName,
		},
		Pods: pod.Functions{
			GetByID:             pod.GetByID(r.Pods),
			GetByPlayerID:       pod.GetByPlayerID(r.Pods),
			Create:              pod.Create(r.Pods, r.PlayerPodRoles),
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
			CreateWithOAuth: user.CreateWithOAuth(r.Users),
		},
	}
}
