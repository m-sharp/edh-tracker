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
	"github.com/m-sharp/edh-tracker/lib/repositories/player"
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

// CreateTestPlayer inserts a fresh player row and returns its ID.
func CreateTestPlayer(t *testing.T, db *gorm.DB) int {
	t.Helper()
	id, err := player.NewRepositoryFromDB(db).Add(context.Background(), fmt.Sprintf("Test Player %d", nextID()))
	require.NoError(t, err)
	return id
}

// CreateTestCommander inserts a fresh commander row and returns its ID.
func CreateTestCommander(t *testing.T, db *gorm.DB) int {
	t.Helper()
	id, err := commander.NewRepositoryFromDB(db).Add(context.Background(), fmt.Sprintf("Test Commander %d", nextID()))
	require.NoError(t, err)
	return id
}

// CreateTestDeck inserts a fresh player + deck row and returns the deck ID.
func CreateTestDeck(t *testing.T, db *gorm.DB) int {
	t.Helper()
	playerID := CreateTestPlayer(t, db)
	id, err := deck.NewRepositoryFromDB(db).Add(context.Background(), playerID, fmt.Sprintf("Test Deck %d", nextID()), 1)
	require.NoError(t, err)
	return id
}
