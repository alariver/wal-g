package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wal-g/wal-g/internal/databases/postgres"
)

func TestTryGetOidPair_DefaultFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/base/1234/5678")
	assert.Equal(t, true, isDB)
	assert.Equal(t, uint32(1234), dbID)
	assert.Equal(t, uint32(5678), tableID)
}

func TestTryGetOidPair_VMFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/base/1234/5678_vm")
	assert.Equal(t, true, isDB)
	assert.Equal(t, uint32(1234), dbID)
	assert.Equal(t, uint32(5678), tableID)
}

func TestTryGetOidPair_PartFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/base/1234/5678.1")
	assert.Equal(t, true, isDB)
	assert.Equal(t, uint32(1234), dbID)
	assert.Equal(t, uint32(5678), tableID)
}

func TestTryGetOidPair_PartFSMFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/base/1234/5678_fsm.1")
	assert.Equal(t, true, isDB)
	assert.Equal(t, uint32(1234), dbID)
	assert.Equal(t, uint32(5678), tableID)
}

func TestTryGetOidPair_TablespaceFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/pg_tblspc/path/in/tablespace/1234/5678")
	assert.Equal(t, true, isDB)
	assert.Equal(t, uint32(1234), dbID)
	assert.Equal(t, uint32(5678), tableID)
}

func TestTryGetOidPair_SpecialFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/base/4/pg_filenode.map")
	assert.Equal(t, true, isDB)
	assert.Equal(t, uint32(4), dbID)
	assert.Equal(t, uint32(0), tableID)
}

func TestTryGetOidPair_RandomFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/path/to/files/1/2")
	assert.Equal(t, false, isDB)
	assert.Equal(t, uint32(0), dbID)
	assert.Equal(t, uint32(0), tableID)
}

func TestTryGetOidPair_DirBetweenFile(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/base/somedir/1/2")
	assert.Equal(t, true, isDB)
	assert.Equal(t, uint32(1), dbID)
	assert.Equal(t, uint32(2), tableID)
}

func TestTryGetOidPair_BaseRoot(t *testing.T) {
	isDB, dbID, tableID := postgres.TryGetOidPair("/base")
	assert.Equal(t, false, isDB)
	assert.Equal(t, uint32(0), dbID)
	assert.Equal(t, uint32(0), tableID)
}

func TestRestoreDesc_OrdinarySkip(t *testing.T) {
	restoreDesc := make(postgres.RestoreDesc)
	restoreDesc.Add(20000, 30000, 100000)
	assert.Equal(t, true, restoreDesc.IsSkipped(20000, 40000))
}

func TestRestoreDesc_NotSpecifiedDatabaseSkip(t *testing.T) {
	restoreDesc := make(postgres.RestoreDesc)
	restoreDesc.Add(20000, 30000, 30000)
	assert.Equal(t, true, restoreDesc.IsSkipped(30000, 10000))
}

func TestRestoreDesc_NoTableSkip(t *testing.T) {
	restoreDesc := make(postgres.RestoreDesc)
	restoreDesc.Add(20000, 30000, 30000)
	assert.Equal(t, false, restoreDesc.IsSkipped(20000, 30000))
}

func TestRestoreDesc_NoSystemSkip(t *testing.T) {
	restoreDesc := make(postgres.RestoreDesc)
	restoreDesc.Add(20000, 30000, 30000)
	assert.Equal(t, true, restoreDesc.IsSkipped(20000, 10000))
}

func TestRestoreDesc_NoDatabaseSkip(t *testing.T) {
	restoreDesc := make(postgres.RestoreDesc)
	restoreDesc.Add(20000, 10000, 10000)
	assert.Equal(t, true, restoreDesc.IsSkipped(20000, 30000))
}

func TestRestoreDesc_NoSystemDatabaseSkip(t *testing.T) {
	restoreDesc := make(postgres.RestoreDesc)
	restoreDesc.Add(20000, 30000, 30000)
	assert.Equal(t, false, restoreDesc.IsSkipped(10000, 40000))
}
