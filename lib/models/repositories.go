package models

import (
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

type Repositories struct {
	Players        *PlayerRepository
	Decks          *DeckRepository
	Games          *GameRepository
	GameResults    *GameResultRepository
	Pods           *PodRepository
	Users          *UserRepository
	Formats        *FormatRepository
	Commanders     *CommanderRepository
	DeckCommanders *DeckCommanderRepository
}

func NewRepositories(log *zap.Logger, client *lib.DBClient) *Repositories {
	gameResults := NewGameResultRepository(client)
	return &Repositories{
		Players:        NewPlayerRepository(client),
		Decks:          NewDeckRepository(client),
		Games:          NewGameRepository(log, client, gameResults),
		GameResults:    gameResults,
		Pods:           NewPodRepository(client),
		Users:          NewUserRepository(client),
		Formats:        NewFormatRepository(client),
		Commanders:     NewCommanderRepository(client),
		DeckCommanders: NewDeckCommanderRepository(client),
	}
}
