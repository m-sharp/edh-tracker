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

var playerColumns = []string{"id", "name", "created_at", "updated_at", "deleted_at"}
var playerStatColumns = []string{"game_id", "place", "kill_count"}
var podIDColumns = []string{"pod_id"}

func TestPlayerRepository_GetById(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(mock sqlmock.Sqlmock)
		wantErr   bool
		checkFunc func(t *testing.T, info *PlayerInfo)
	}{
		{
			name: "Success",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerByID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerColumns).
						AddRow(1, "Alice", time.Time{}, time.Time{}, nil))
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerStats)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerStatColumns).
						AddRow(100, 1, 2))
				mock.ExpectQuery(regexp.QuoteMeta(GetPodIDsByPlayerID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(podIDColumns).
						AddRow(1).AddRow(3))
			},
			wantErr: false,
			checkFunc: func(t *testing.T, info *PlayerInfo) {
				assert.Equal(t, "Alice", info.Name)
				assert.Equal(t, 5, info.Stats.Points) // kills=2 + place=1 bonus=3
				assert.Equal(t, []int{1, 3}, info.PodIDs)
			},
		},
		{
			name: "NoGames_NoPods",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerByID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerColumns).
						AddRow(1, "Bob", time.Time{}, time.Time{}, nil))
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerStats)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerStatColumns))
				mock.ExpectQuery(regexp.QuoteMeta(GetPodIDsByPlayerID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(podIDColumns))
			},
			wantErr: false,
			checkFunc: func(t *testing.T, info *PlayerInfo) {
				assert.Equal(t, 0, info.Stats.Games)
				assert.Equal(t, 0, info.Stats.Points)
				assert.Equal(t, []int{}, info.PodIDs)
			},
		},
		{
			name: "NotFound",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerByID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerColumns))
			},
			wantErr: true,
		},
		{
			name: "TooMany",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerByID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerColumns).
						AddRow(1, "Alice", time.Time{}, time.Time{}, nil).
						AddRow(2, "Bob", time.Time{}, time.Time{}, nil))
			},
			wantErr: true,
		},
		{
			name: "StatsError",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerByID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerColumns).
						AddRow(1, "Alice", time.Time{}, time.Time{}, nil))
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerStats)).
					WithArgs(1).
					WillReturnError(errors.New("stats error"))
			},
			wantErr: true,
		},
		{
			name: "PodsError",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerByID)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerColumns).
						AddRow(1, "Alice", time.Time{}, time.Time{}, nil))
				mock.ExpectQuery(regexp.QuoteMeta(GetPlayerStats)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(playerStatColumns))
				mock.ExpectQuery(regexp.QuoteMeta(GetPodIDsByPlayerID)).
					WithArgs(1).
					WillReturnError(errors.New("pods error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			client, mock := newMockDB(tt)
			tc.setup(mock)

			repo := NewPlayerRepository(client)
			info, err := repo.GetById(ctx, 1)

			if tc.wantErr {
				assert.Error(tt, err)
			} else {
				require.NoError(tt, err)
				require.NotNil(tt, info)
				tc.checkFunc(tt, info)
			}
			assert.NoError(tt, mock.ExpectationsWereMet())
		})
	}
}
