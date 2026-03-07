package playerPodRole

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

func TestGetRole_Found(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "pod_id", "player_id", "role", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, 10, 20, "manager", now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getRole)).WithArgs(10, 20).WillReturnRows(rows)

	got, err := repo.GetRole(context.Background(), 10, 20)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "manager", got.Role)
	assert.Equal(t, 10, got.PodID)
	assert.Equal(t, 20, got.PlayerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRole_NotFound(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"id", "pod_id", "player_id", "role", "created_at", "updated_at", "deleted_at"})
	mock.ExpectQuery(regexp.QuoteMeta(getRole)).WithArgs(10, 20).WillReturnRows(rows)

	got, err := repo.GetRole(context.Background(), 10, 20)
	require.NoError(t, err)
	assert.Nil(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetRole_Success(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	mock.ExpectExec(regexp.QuoteMeta(setRole)).
		WithArgs(10, 20, "manager").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.SetRole(context.Background(), 10, 20, "manager")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMembersWithRoles_Empty(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	rows := sqlmock.NewRows([]string{"id", "pod_id", "player_id", "role", "created_at", "updated_at", "deleted_at"})
	mock.ExpectQuery(regexp.QuoteMeta(getMembersWithRole)).WithArgs(10).WillReturnRows(rows)

	got, err := repo.GetMembersWithRoles(context.Background(), 10)
	require.NoError(t, err)
	assert.Empty(t, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMembersWithRoles_Rows(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "pod_id", "player_id", "role", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, 10, 1, "manager", now, now, nil).
		AddRow(2, 10, 2, "member", now, now, nil)
	mock.ExpectQuery(regexp.QuoteMeta(getMembersWithRole)).WithArgs(10).WillReturnRows(rows)

	got, err := repo.GetMembersWithRoles(context.Background(), 10)
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "manager", got[0].Role)
	assert.Equal(t, "member", got[1].Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBulkAdd_EmptyPlayerIDs(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	// No DB calls expected for empty input
	err := repo.BulkAdd(context.Background(), 10, []int{}, "member")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBulkAdd_MultiplePlayerIDs(t *testing.T) {
	client, mock := newMockDB(t)
	repo := NewRepository(client)

	expectedQuery := "INSERT INTO player_pod_role (pod_id, player_id, role) VALUES (?,?,?),(?,?,?)"
	mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
		WithArgs(10, 1, "member", 10, 2, "member").
		WillReturnResult(sqlmock.NewResult(2, 2))

	err := repo.BulkAdd(context.Background(), 10, []int{1, 2}, "member")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
