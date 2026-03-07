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

func TestUpdate_NameOnly(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	name := "New Name"
	mock.ExpectExec(`UPDATE deck SET name = \? WHERE id = \? AND deleted_at IS NULL`).
		WithArgs("New Name", 7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), 7, UpdateFields{Name: &name})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_MultipleFields(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	name := "Renamed"
	formatID := 2
	retired := true
	mock.ExpectExec(`UPDATE deck SET`).
		WithArgs("Renamed", 2, true, 7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), 7, UpdateFields{Name: &name, FormatID: &formatID, Retired: &retired})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NoFields_NoOp(t *testing.T) {
	client, _ := newMockDB(t)
	repo := NewRepository(client)

	// No mock expectations — no query should be issued
	err := repo.Update(context.Background(), 7, UpdateFields{})
	require.NoError(t, err)
}

func TestGetAllByPlayerIDs_Empty(t *testing.T) {
	client, _ := newMockDB(t)
	repo := NewRepository(client)

	got, err := repo.GetAllByPlayerIDs(context.Background(), []int{})
	require.NoError(t, err)
	assert.Len(t, got, 0)
}
