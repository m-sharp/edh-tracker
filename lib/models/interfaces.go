package models

import "context"

type PlayerRepositoryInterface interface {
	GetAll(ctx context.Context) ([]PlayerInfo, error)
	GetById(ctx context.Context, playerID int) (*PlayerInfo, error)
	Add(ctx context.Context, name string) (int, error)
}

type DeckRepositoryInterface interface {
	GetAll(ctx context.Context) ([]DeckWithStats, error)
	GetAllForPlayer(ctx context.Context, playerID int) ([]DeckWithStats, error)
	GetById(ctx context.Context, deckID int) (*DeckWithStats, error)
	Add(ctx context.Context, playerID int, name string, formatID int) (int, error)
	Retire(ctx context.Context, deckID int) error
}

type GameRepositoryInterface interface {
	GetAllByPod(ctx context.Context, podId int) ([]GameDetails, error)
	GetAllByDeck(ctx context.Context, deckId int) ([]GameDetails, error)
	GetGameById(ctx context.Context, gameId int) (*GameDetails, error)
	Add(ctx context.Context, description string, podID, formatID int, results ...GameResult) error
}

type FormatRepositoryInterface interface {
	GetAll(ctx context.Context) ([]Format, error)
	GetById(ctx context.Context, id int) (*Format, error)
}

type DeckCommanderRepositoryInterface interface {
	Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error)
}

type CommanderRepositoryInterface interface {
	GetById(ctx context.Context, id int) (*Commander, error)
	Add(ctx context.Context, name string) (int, error)
}

type PodRepositoryInterface interface {
	GetByID(ctx context.Context, podID int) (*Pod, error)
	GetByPlayerID(ctx context.Context, playerID int) ([]Pod, error)
	Add(ctx context.Context, name string) (int, error)
	AddPlayerToPod(ctx context.Context, podID, playerID int) error
}
