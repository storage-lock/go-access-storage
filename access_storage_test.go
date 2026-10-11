package access_storage

import (
	"context"
	"database/sql"
	"testing"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewAccessStorage_Valid(t *testing.T) {
	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	s, err := NewAccessStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, StorageName, s.GetName())
}

func TestAccessStorage_GetName(t *testing.T) {
	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	s, err := NewAccessStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, StorageName, s.GetName())
}

func TestAccessStorage_Capabilities(t *testing.T) {
	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	s, err := NewAccessStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	caps := s.Capabilities()
	assert.Contains(t, caps, storage.CapabilityCAS)
	assert.Contains(t, caps, storage.CapabilityReliableTime)
}

func TestAccessStorage_Close(t *testing.T) {
	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{}))
	s, err := NewAccessStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	err = s.Close(context.Background())
	assert.NoError(t, err)
}
