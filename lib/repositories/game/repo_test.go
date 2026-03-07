package game

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

func gameColumns() []string {
	return []string{"id", "description", "pod_id", "format_id", "created_at", "updated_at", "deleted_at"}
}

func TestGetById_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	now := time.Now()
	rows := sqlmock.NewRows(gameColumns()).
		AddRow(5, "Friday Night", 2, 1, now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getGameByID)).WithArgs(5).WillReturnRows(rows)

	got, err := repo.GetById(context.Background(), 5)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 5, got.ID)
	assert.Equal(t, "Friday Night", got.Description)
	assert.Equal(t, 2, got.PodID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetById_NotFound(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows(gameColumns())
	mock.ExpectQuery(regexp.QuoteMeta(getGameByID)).WithArgs(99).WillReturnRows(rows)

	got, err := repo.GetById(context.Background(), 99)
	require.NoError(t, err)
	assert.Nil(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(updateGame)).
		WithArgs("New Description", 5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), 5, "New Description")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NotFound(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(updateGame)).
		WithArgs("New Description", 99).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(context.Background(), 99, "New Description")
	assert.ErrorContains(t, err, "unexpected number of rows")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllByPlayerID_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	now := time.Now()
	rows := sqlmock.NewRows(gameColumns()).
		AddRow(1, "Game 1", 2, 1, now, now, nil).
		AddRow(2, "Game 2", 2, 1, now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getGamesByPlayerID)).WithArgs(7).WillReturnRows(rows)

	got, err := repo.GetAllByPlayerID(context.Background(), 7)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllByPlayerID_Empty(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows(gameColumns())
	mock.ExpectQuery(regexp.QuoteMeta(getGamesByPlayerID)).WithArgs(99).WillReturnRows(rows)

	got, err := repo.GetAllByPlayerID(context.Background(), 99)
	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSoftDelete_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(softDeleteGame)).
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.SoftDelete(context.Background(), 5)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
