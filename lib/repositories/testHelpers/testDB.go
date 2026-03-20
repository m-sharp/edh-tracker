package testHelpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewTestDB stands up a connection to the local DB for integration testing, wrapping all calls into a transaction that
// is rolled back at the end of testing.
// Need to have docker db running locally to run int tests.
// Yeah, I've checked in my fake root password...o noes don't hack my local test db pl0x.
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"root",
		"BirdSquad",
		"host.docker.internal",
		"3306",
		"pod_tracker",
	)
	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	tx := db.Begin()
	require.NoError(t, tx.Error)
	t.Cleanup(func() { tx.Rollback() })

	return tx
}
