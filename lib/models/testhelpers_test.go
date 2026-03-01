package models

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/m-sharp/edh-tracker/lib"
	"github.com/stretchr/testify/require"
)

func newMockDB(t *testing.T) (*lib.DBClient, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	return &lib.DBClient{Db: sqlx.NewDb(db, "sqlmock")}, mock
}
