package models

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gameResultColumns = []string{
	"id", "game_id", "deck_id", "deck_name",
	"commander_name", "partner_commander_name",
	"place", "kill_count",
	"created_at", "updated_at", "deleted_at",
}

func TestGameResultRepository_GetByGameId(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(mock sqlmock.Sqlmock)
		wantErr   bool
		checkFunc func(t *testing.T, results []GameResult)
	}{
		{
			name: "WithCommanders",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(gameResultColumns).
					AddRow(1, 10, 5, "My Deck", "Atraxa", "Thrasios", 1, 2, time.Time{}, time.Time{}, nil)
				mock.ExpectQuery(regexp.QuoteMeta(GetGameResultsByGameID)).
					WithArgs(10).
					WillReturnRows(rows)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, results []GameResult) {
				require.Len(t, results, 1)
				r := results[0]
				assert.Equal(t, 5, r.Points) // kills=2 + place=1 bonus=3
				assert.NotNil(t, r.CommanderName)
				assert.Equal(t, "Atraxa", *r.CommanderName)
				assert.NotNil(t, r.PartnerCommanderName)
				assert.Equal(t, "Thrasios", *r.PartnerCommanderName)
			},
		},
		{
			name: "NoCommanders",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(gameResultColumns).
					AddRow(2, 10, 5, "My Deck", nil, nil, 2, 1, time.Time{}, time.Time{}, nil)
				mock.ExpectQuery(regexp.QuoteMeta(GetGameResultsByGameID)).
					WithArgs(10).
					WillReturnRows(rows)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, results []GameResult) {
				require.Len(t, results, 1)
				r := results[0]
				assert.Equal(t, 3, r.Points) // kills=1 + place=2 bonus=2
				assert.Nil(t, r.CommanderName)
				assert.Nil(t, r.PartnerCommanderName)
			},
		},
		{
			name: "SoloCommander",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(gameResultColumns).
					AddRow(3, 10, 5, "My Deck", "Kenrith", nil, 3, 0, time.Time{}, time.Time{}, nil)
				mock.ExpectQuery(regexp.QuoteMeta(GetGameResultsByGameID)).
					WithArgs(10).
					WillReturnRows(rows)
			},
			wantErr: false,
			checkFunc: func(t *testing.T, results []GameResult) {
				require.Len(t, results, 1)
				r := results[0]
				assert.NotNil(t, r.CommanderName)
				assert.Equal(t, "Kenrith", *r.CommanderName)
				assert.Nil(t, r.PartnerCommanderName)
			},
		},
		{
			name: "Empty",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(gameResultColumns)
				mock.ExpectQuery(regexp.QuoteMeta(GetGameResultsByGameID)).
					WithArgs(10).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "DBError",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetGameResultsByGameID)).
					WithArgs(10).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			client, mock := newMockDB(tt)
			tc.setup(mock)

			repo := NewGameResultRepository(client)
			results, err := repo.GetByGameId(ctx, 10)

			if tc.wantErr {
				assert.Error(tt, err)
			} else {
				require.NoError(tt, err)
				tc.checkFunc(tt, results)
			}
			assert.NoError(tt, mock.ExpectationsWereMet())
		})
	}
}

func TestGameResultRepository_SoftDelete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(SoftDeleteGameResult)).
					WithArgs(42).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "ZeroRows",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(SoftDeleteGameResult)).
					WithArgs(42).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "ExecError",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(SoftDeleteGameResult)).
					WithArgs(42).
					WillReturnError(errors.New("exec error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			client, mock := newMockDB(tt)
			tc.setup(mock)

			repo := NewGameResultRepository(client)
			err := repo.SoftDelete(ctx, 42)

			if tc.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.NoError(tt, mock.ExpectationsWereMet())
		})
	}
}
