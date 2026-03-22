package testHelpers

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

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

var fixtureCounter int64

func nextID() int64 {
	return atomic.AddInt64(&fixtureCounter, 1)
}

func NewPlayerRepo(db *gorm.DB) *player.Repository {
	return player.NewRepositoryFromDB(db)
}

func NewFormatRepo(db *gorm.DB) *format.Repository {
	return format.NewRepositoryFromDB(db)
}

func NewCommanderRepo(db *gorm.DB) *commander.Repository {
	return commander.NewRepositoryFromDB(db)
}

func NewDeckRepo(db *gorm.DB) *deck.Repository {
	return deck.NewRepositoryFromDB(db)
}

func NewDeckCommanderRepo(db *gorm.DB) *deckCommander.Repository {
	return deckCommander.NewRepositoryFromDB(db)
}

func NewPodRepo(db *gorm.DB) *pod.Repository {
	return pod.NewRepositoryFromDB(db)
}

func NewPlayerPodRoleRepo(db *gorm.DB) *playerPodRole.Repository {
	return playerPodRole.NewRepositoryFromDB(db)
}

func NewGameRepo(db *gorm.DB) *game.Repository {
	return game.NewRepositoryFromDB(db)
}

func NewGameResultRepo(db *gorm.DB) *gameResult.Repository {
	return gameResult.NewRepositoryFromDB(db)
}

func NewPodInviteRepo(db *gorm.DB) *podInvite.Repository {
	return podInvite.NewRepositoryFromDB(db)
}

func NewUserRepo(db *gorm.DB) *user.Repository {
	return user.NewRepositoryFromDB(db)
}

// CreateTestGameResult inserts a fresh game_result row and returns its ID.
func CreateTestGameResult(t *testing.T, db *gorm.DB, gameID, deckID, place, killCount int) int {
	t.Helper()
	repo := NewGameResultRepo(db)
	id, err := repo.Add(context.Background(), gameResult.Model{
		GameID:    gameID,
		DeckID:    deckID,
		Place:     place,
		KillCount: killCount,
	})
	require.NoError(t, err)
	return id
}

// CreateTestPod inserts a fresh pod row and returns its ID.
func CreateTestPod(t *testing.T, db *gorm.DB) int {
	t.Helper()
	repo := NewPodRepo(db)
	id, err := repo.Add(context.Background(), fmt.Sprintf("TestPod-%d", nextID()))
	require.NoError(t, err)
	return id
}

// CreateTestPlayer inserts a fresh player row and returns its ID.
func CreateTestPlayer(t *testing.T, db *gorm.DB) int {
	t.Helper()
	repo := NewPlayerRepo(db)
	id, err := repo.Add(context.Background(), fmt.Sprintf("Test Player %d", nextID()))
	require.NoError(t, err)
	return id
}

// CreateTestCommander inserts a fresh commander row and returns its ID.
func CreateTestCommander(t *testing.T, db *gorm.DB) int {
	t.Helper()
	repo := NewCommanderRepo(db)
	id, err := repo.Add(context.Background(), fmt.Sprintf("Test Commander %d", nextID()))
	require.NoError(t, err)
	return id
}

// GetCommanderFormatID looks up the ID of the "commander" format in the DB.
func GetCommanderFormatID(t *testing.T, db *gorm.DB) int {
	t.Helper()
	repo := NewFormatRepo(db)
	f, err := repo.GetByName(context.Background(), "commander")
	require.NoError(t, err)
	require.NotNil(t, f, "commander format not found in DB")
	return f.ID
}

// CreateTestGame inserts a fresh pod + game row and returns the game ID.
func CreateTestGame(t *testing.T, db *gorm.DB) int {
	t.Helper()
	podID := CreateTestPod(t, db)
	formatID := GetCommanderFormatID(t, db)

	repo := NewGameRepo(db)
	id, err := repo.Add(context.Background(), fmt.Sprintf("Test Game %d", nextID()), podID, formatID)
	require.NoError(t, err)
	return id
}

// CreateTestDeck inserts a fresh player + deck row and returns the deck Model.
func CreateTestDeck(t *testing.T, db *gorm.DB) deck.Model {
	t.Helper()
	playerID := CreateTestPlayer(t, db)
	name := fmt.Sprintf("Test Deck %d", nextID())
	formatID := GetCommanderFormatID(t, db)

	repo := NewDeckRepo(db)
	id, err := repo.Add(context.Background(), playerID, name, formatID)
	require.NoError(t, err)

	createdDeck, err := repo.GetById(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, createdDeck)

	return *createdDeck
}

// CreateTestDeckWithCommander inserts a fresh player + deck + deck_commander row and returns the deck Model.
func CreateTestDeckWithCommander(t *testing.T, db *gorm.DB) deck.Model {
	t.Helper()
	testDeck := CreateTestDeck(t, db)
	commanderID := CreateTestCommander(t, db)

	dcRepo := NewDeckCommanderRepo(db)
	_, err := dcRepo.Add(context.Background(), testDeck.ID, commanderID, nil)
	require.NoError(t, err)

	return testDeck
}
