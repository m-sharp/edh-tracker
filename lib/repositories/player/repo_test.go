package player

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
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(42, "Alice", now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getPlayerByID)).WithArgs(42).WillReturnRows(rows)

	got, err := repo.GetById(context.Background(), 42)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 42, got.ID)
	assert.Equal(t, "Alice", got.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetById_NotFound(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"})
	mock.ExpectQuery(regexp.QuoteMeta(getPlayerByID)).WithArgs(99).WillReturnRows(rows)

	got, err := repo.GetById(context.Background(), 99)
	require.NoError(t, err)
	assert.Nil(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSoftDelete_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(softDeletePlayer)).
		WithArgs(42).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SoftDelete(context.Background(), 42)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(updatePlayer)).
		WithArgs("NewName", 42).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), 42, "NewName")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NotFound(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(updatePlayer)).
		WithArgs("NewName", 99).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(context.Background(), 99, "NewName")
	assert.ErrorContains(t, err, "unexpected number of rows")
	assert.NoError(t, mock.ExpectationsWereMet())
}
