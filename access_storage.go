package access_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
)

// AccessStorage Microsoft Access 数据库存储实现
type AccessStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewAccessStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*AccessStorage, error) {
	baseStorage, err := sql_based_storage.NewSqlBasedStorage(options)
	if err != nil {
		return nil, err
	}
	return &AccessStorage{SqlBasedStorage: baseStorage}, nil
}

const StorageName = "access-storage"

func (x *AccessStorage) GetName() string { return StorageName }

func (x *AccessStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{storage.CapabilityCAS, storage.CapabilityReliableTime}
}
