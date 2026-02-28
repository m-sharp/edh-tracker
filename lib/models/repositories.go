package models

import (
	"go.uber.org/zap"

	"github.com/m-sharp/edh-tracker/lib"
)

type Repositories struct {
	Players *PlayerProvider
	Decks   *DeckProvider
	Games   *GameProvider
	Pods    *PodProvider
	Users   *UserProvider
}

func NewRepositories(log *zap.Logger, client *lib.DBClient) *Repositories {
	return &Repositories{
		Players: NewPlayerProvider(client),
		Decks:   NewDeckProvider(client),
		Games:   NewGameProvider(log, client),
		Pods:    NewPodProvider(client),
		Users:   NewUserProvider(client),
	}
}
