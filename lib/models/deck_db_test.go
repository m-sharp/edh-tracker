package models

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var deckColumns = []string{
	"id", "player_id", "player_name", "name", "format_id", "format_name", "retired",
	"created_at", "updated_at", "deleted_at",
	"commander_id", "commander_name", "partner_commander_id", "partner_commander_name",
}

var deckStatColumns = []string{"game_id", "place", "kill_count"}

// deckBaseRow returns a standard deck row with nulled commander fields.
func deckBaseRow(commanderID sql.NullInt64, commanderName sql.NullString, partnerID sql.NullInt64, partnerName sql.NullString) *sqlmock.Rows {
	return sqlmock.NewRows(deckColumns).AddRow(
		1, 2, "Alice", "Test Deck", 1, "commander", false,
		time.Time{}, time.Time{}, nil,
		commanderID, commanderName, partnerID, partnerName,
	)
}

func TestDeckRepository_GetById(t *testing.T) {
	ctx := context.Background()

	nullInt := sql.NullInt64{Valid: false}
	nullStr := sql.NullString{Valid: false}
	cmdInt := sql.NullInt64{Valid: true, Int64: 1}
	cmdStr := sql.NullString{Valid: true, String: "Atraxa"}
	partInt := sql.NullInt64{Valid: true, Int64: 2}
	partStr := sql.NullString{Valid: true, String: "Thrasios"}

	tests := []struct {
		name      string
		setup     func(mock sqlmock.Sqlmock)
		wantErr   bool
		checkFunc func(t *testing.T, deck *DeckWithStats)
	}{
		{
			name: "NoCommander",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckByID)).
					WithArgs(1).
					WillReturnRows(deckBaseRow(nullInt, nullStr, nullInt, nullStr))
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckStats)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(deckStatColumns))
			},
			wantErr: false,
			checkFunc: func(t *testing.T, deck *DeckWithStats) {
				assert.Nil(t, deck.Commanders)
			},
		},
		{
			name: "SoloCommander",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckByID)).
					WithArgs(1).
					WillReturnRows(deckBaseRow(cmdInt, cmdStr, nullInt, nullStr))
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckStats)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(deckStatColumns))
			},
			wantErr: false,
			checkFunc: func(t *testing.T, deck *DeckWithStats) {
				require.NotNil(t, deck.Commanders)
				assert.Equal(t, 1, deck.Commanders.CommanderID)
				assert.Equal(t, "Atraxa", deck.Commanders.CommanderName)
				assert.Nil(t, deck.Commanders.PartnerCommanderID)
				assert.Nil(t, deck.Commanders.PartnerCommanderName)
			},
		},
		{
			name: "PartnerCommander",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckByID)).
					WithArgs(1).
					WillReturnRows(deckBaseRow(cmdInt, cmdStr, partInt, partStr))
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckStats)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(deckStatColumns))
			},
			wantErr: false,
			checkFunc: func(t *testing.T, deck *DeckWithStats) {
				require.NotNil(t, deck.Commanders)
				assert.Equal(t, 1, deck.Commanders.CommanderID)
				require.NotNil(t, deck.Commanders.PartnerCommanderID)
				assert.Equal(t, 2, *deck.Commanders.PartnerCommanderID)
				require.NotNil(t, deck.Commanders.PartnerCommanderName)
				assert.Equal(t, "Thrasios", *deck.Commanders.PartnerCommanderName)
			},
		},
		{
			name: "WithStats",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckByID)).
					WithArgs(1).
					WillReturnRows(deckBaseRow(nullInt, nullStr, nullInt, nullStr))
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckStats)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(deckStatColumns).AddRow(100, 1, 2))
			},
			wantErr: false,
			checkFunc: func(t *testing.T, deck *DeckWithStats) {
				assert.Equal(t, 1, deck.Stats.Games)
				assert.Equal(t, 2, deck.Stats.Kills)
				assert.Equal(t, 5, deck.Stats.Points) // kills=2 + place=1 bonus=3
			},
		},
		{
			name: "NotFound",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckByID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(deckColumns))
			},
			wantErr: true,
		},
		{
			name: "DBError",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetDeckByID)).
					WithArgs(1).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			client, mock := newMockDB(tt)
			tc.setup(mock)

			repo := NewDeckRepository(client)
			deck, err := repo.GetById(ctx, 1)

			if tc.wantErr {
				assert.Error(tt, err)
			} else {
				require.NoError(tt, err)
				require.NotNil(tt, deck)
				tc.checkFunc(tt, deck)
			}
			assert.NoError(tt, mock.ExpectationsWereMet())
		})
	}
}

func TestDeckRepository_SoftDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(SoftDeleteDeck)).
					WithArgs(7).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "ZeroRows",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(SoftDeleteDeck)).
					WithArgs(7).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			client, mock := newMockDB(tt)
			tc.setup(mock)

			repo := NewDeckRepository(client)
			err := repo.SoftDelete(ctx, 7)

			if tc.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.NoError(tt, mock.ExpectationsWereMet())
		})
	}
}
