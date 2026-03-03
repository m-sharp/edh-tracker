package deck

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

func TestGetById_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "player_id", "name", "format_id", "retired", "created_at", "updated_at", "deleted_at"}).
		AddRow(7, 3, "Krenko Goblins", 1, false, now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getDeckByID)).WithArgs(7).WillReturnRows(rows)

	got, err := repo.GetById(context.Background(), 7)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 7, got.ID)
	assert.Equal(t, 3, got.PlayerID)
	assert.Equal(t, "Krenko Goblins", got.Name)
	assert.Equal(t, 1, got.FormatID)
	assert.False(t, got.Retired)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetById_NotFound(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"id", "player_id", "name", "format_id", "retired", "created_at", "updated_at", "deleted_at"})
	mock.ExpectQuery(regexp.QuoteMeta(getDeckByID)).WithArgs(99).WillReturnRows(rows)

	got, err := repo.GetById(context.Background(), 99)
	require.NoError(t, err)
	assert.Nil(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSoftDelete_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(softDeleteDeck)).
		WithArgs(7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SoftDelete(context.Background(), 7)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
