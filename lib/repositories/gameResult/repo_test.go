package gameResult

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m-sharp/edh-tracker/lib"
)

func newMockDB(t *testing.T) (*lib.DBClient, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return &lib.DBClient{Db: sqlx.NewDb(db, "sqlmock")}, mock
}

func TestGetByGameId_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "game_id", "deck_id", "place", "kill_count", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, 5, 10, 1, 2, now, now, nil).
		AddRow(2, 5, 11, 2, 0, now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getGameResultsByGameID)).WithArgs(5).WillReturnRows(rows)

	got, err := repo.GetByGameId(context.Background(), 5)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 1, got[0].ID)
	assert.Equal(t, 5, got[0].GameID)
	assert.Equal(t, 10, got[0].DeckID)
	assert.Equal(t, 1, got[0].Place)
	assert.Equal(t, 2, got[0].KillCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByGameId_Empty(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"id", "game_id", "deck_id", "place", "kill_count", "created_at", "updated_at", "deleted_at"})
	mock.ExpectQuery(regexp.QuoteMeta(getGameResultsByGameID)).WithArgs(99).WillReturnRows(rows)

	got, err := repo.GetByGameId(context.Background(), 99)
	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGameResultSoftDelete_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(softDeleteGameResult)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SoftDelete(context.Background(), 1)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStatsForPlayer_WithGames(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"game_id", "place", "kill_count", "player_count"}).
		AddRow(1, 1, 2, 4).
		AddRow(2, 2, 0, 4).
		AddRow(3, 1, 1, 4)
	mock.ExpectQuery(regexp.QuoteMeta(getStatsForPlayer)).WithArgs(7).WillReturnRows(rows)

	got, err := repo.GetStatsForPlayer(context.Background(), 7)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 3, got.Games)
	assert.Equal(t, 3, got.Kills)
	assert.Equal(t, 11, got.Points) // 5 + 2 + 4
	assert.Equal(t, map[int]int{1: 2, 2: 1}, got.Record)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStatsForDeck_Empty(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"game_id", "place", "kill_count", "player_count"})
	mock.ExpectQuery(regexp.QuoteMeta(getStatsForDeck)).WithArgs(10).WillReturnRows(rows)

	got, err := repo.GetStatsForDeck(context.Background(), 10)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 0, got.Games)
	assert.Len(t, got.Record, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}
