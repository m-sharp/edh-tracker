package routers

import (
	"context"

	"github.com/m-sharp/edh-tracker/lib/models"
)

// mockPlayerRepo implements models.PlayerRepositoryInterface for testing.
type mockPlayerRepo struct {
	GetAllFn  func(ctx context.Context) ([]models.PlayerInfo, error)
	GetByIdFn func(ctx context.Context, id int) (*models.PlayerInfo, error)
	AddFn     func(ctx context.Context, name string) (int, error)
}

func (m *mockPlayerRepo) GetAll(ctx context.Context) ([]models.PlayerInfo, error) {
	return m.GetAllFn(ctx)
}
func (m *mockPlayerRepo) GetById(ctx context.Context, id int) (*models.PlayerInfo, error) {
	return m.GetByIdFn(ctx, id)
}
func (m *mockPlayerRepo) Add(ctx context.Context, name string) (int, error) {
	return m.AddFn(ctx, name)
}

// mockDeckRepo implements models.DeckRepositoryInterface for testing.
type mockDeckRepo struct {
	GetAllFn          func(ctx context.Context) ([]models.DeckWithStats, error)
	GetAllForPlayerFn func(ctx context.Context, playerID int) ([]models.DeckWithStats, error)
	GetByIdFn         func(ctx context.Context, deckID int) (*models.DeckWithStats, error)
	AddFn             func(ctx context.Context, playerID int, name string, formatID int) (int, error)
	RetireFn          func(ctx context.Context, deckID int) error
}

func (m *mockDeckRepo) GetAll(ctx context.Context) ([]models.DeckWithStats, error) {
	return m.GetAllFn(ctx)
}
func (m *mockDeckRepo) GetAllForPlayer(ctx context.Context, playerID int) ([]models.DeckWithStats, error) {
	return m.GetAllForPlayerFn(ctx, playerID)
}
func (m *mockDeckRepo) GetById(ctx context.Context, deckID int) (*models.DeckWithStats, error) {
	return m.GetByIdFn(ctx, deckID)
}
func (m *mockDeckRepo) Add(ctx context.Context, playerID int, name string, formatID int) (int, error) {
	return m.AddFn(ctx, playerID, name, formatID)
}
func (m *mockDeckRepo) Retire(ctx context.Context, deckID int) error {
	return m.RetireFn(ctx, deckID)
}

// mockDeckCommanderRepo implements models.DeckCommanderRepositoryInterface for testing.
type mockDeckCommanderRepo struct {
	AddFn func(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error)
}

func (m *mockDeckCommanderRepo) Add(ctx context.Context, deckID, commanderID int, partnerCommanderID *int) (int, error) {
	return m.AddFn(ctx, deckID, commanderID, partnerCommanderID)
}

// mockFormatRepo implements models.FormatRepositoryInterface for testing.
type mockFormatRepo struct {
	GetAllFn  func(ctx context.Context) ([]models.Format, error)
	GetByIdFn func(ctx context.Context, id int) (*models.Format, error)
}

func (m *mockFormatRepo) GetAll(ctx context.Context) ([]models.Format, error) {
	return m.GetAllFn(ctx)
}
func (m *mockFormatRepo) GetById(ctx context.Context, id int) (*models.Format, error) {
	return m.GetByIdFn(ctx, id)
}

// mockGameRepo implements models.GameRepositoryInterface for testing.
type mockGameRepo struct {
	GetAllByPodFn  func(ctx context.Context, podId int) ([]models.GameDetails, error)
	GetAllByDeckFn func(ctx context.Context, deckId int) ([]models.GameDetails, error)
	GetGameByIdFn  func(ctx context.Context, gameId int) (*models.GameDetails, error)
	AddFn          func(ctx context.Context, description string, podID, formatID int, results ...models.GameResult) error
}

func (m *mockGameRepo) GetAllByPod(ctx context.Context, podId int) ([]models.GameDetails, error) {
	return m.GetAllByPodFn(ctx, podId)
}
func (m *mockGameRepo) GetAllByDeck(ctx context.Context, deckId int) ([]models.GameDetails, error) {
	return m.GetAllByDeckFn(ctx, deckId)
}
func (m *mockGameRepo) GetGameById(ctx context.Context, gameId int) (*models.GameDetails, error) {
	return m.GetGameByIdFn(ctx, gameId)
}
func (m *mockGameRepo) Add(ctx context.Context, description string, podID, formatID int, results ...models.GameResult) error {
	return m.AddFn(ctx, description, podID, formatID, results...)
}
