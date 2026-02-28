package models

import (
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

type Repositories struct {
	Players        *PlayerProvider
	Decks          *DeckProvider
	Games          *GameProvider
	GameResults    *GameResultProvider
	Pods           *PodProvider
	Users          *UserProvider
	Formats        *FormatProvider
	Commanders     *CommanderProvider
	DeckCommanders *DeckCommanderProvider
}

func NewRepositories(log *zap.Logger, client *lib.DBClient) *Repositories {
	gameResults := NewGameResultProvider(client)
	return &Repositories{
		Players:        NewPlayerProvider(client),
		Decks:          NewDeckProvider(client),
		Games:          NewGameProvider(log, client, gameResults),
		GameResults:    gameResults,
		Pods:           NewPodProvider(client),
		Users:          NewUserProvider(client),
		Formats:        NewFormatProvider(client),
		Commanders:     NewCommanderProvider(client),
		DeckCommanders: NewDeckCommanderProvider(client),
	}
}
