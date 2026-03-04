package podInvite

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

func TestGetByCode_Found(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	now := time.Now()
	future := now.Add(7 * 24 * time.Hour)
	rows := sqlmock.NewRows([]string{"id", "pod_id", "invite_code", "created_by_player_id", "expires_at", "used_count", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, 5, "abc-123", 10, future, 0, now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getByCode)).WithArgs("abc-123").WillReturnRows(rows)

	got, err := repo.GetByCode(context.Background(), "abc-123")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 5, got.PodID)
	assert.Equal(t, "abc-123", got.InviteCode)
	assert.Equal(t, 10, got.CreatedByPlayerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByCode_NotFound(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"id", "pod_id", "invite_code", "created_by_player_id", "expires_at", "used_count", "created_at", "updated_at", "deleted_at"})
	mock.ExpectQuery(regexp.QuoteMeta(getByCode)).WithArgs("missing").WillReturnRows(rows)

	got, err := repo.GetByCode(context.Background(), "missing")
	require.NoError(t, err)
	assert.Nil(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdd_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	future := time.Now().Add(7 * 24 * time.Hour)
	mock.ExpectExec(regexp.QuoteMeta(insertInvite)).
		WithArgs(5, "abc-123", 10, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Add(context.Background(), 5, 10, "abc-123", &future)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIncrementUsedCount_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(incrementUsedCount)).
		WithArgs("abc-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.IncrementUsedCount(context.Background(), "abc-123")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
