package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDataDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	os.Setenv("DATA_DIR", dir)
	t.Cleanup(func() { os.Unsetenv("DATA_DIR") })
	return dir
}

func TestMigration013_FreshDB(t *testing.T) {
	dataDir := testDataDir(t)
	db, err := Open(dataDir)
	require.NoError(t, err)
	defer db.Close()

	exists, err := db.columnExists("users", "pin_hash")
	require.NoError(t, err)
	assert.True(t, exists)

	var version int
	require.NoError(t, db.QueryRow("SELECT MAX(version) FROM schema_version").Scan(&version))
	assert.Equal(t, 20, version)
}

func TestMigration013_PartiallyMigratedDB(t *testing.T) {
	dataDir := testDataDir(t)
	db, err := Open(dataDir)
	require.NoError(t, err)

	_, err = db.Exec("DELETE FROM schema_version WHERE version = 13")
	require.NoError(t, err)
	require.NoError(t, db.Close())

	db2, err := Open(dataDir)
	require.NoError(t, err)
	defer db2.Close()

	var version int
	require.NoError(t, db2.QueryRow("SELECT MAX(version) FROM schema_version").Scan(&version))
	assert.Equal(t, 20, version)
}
