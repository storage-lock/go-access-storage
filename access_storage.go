package access_storage

import (
	"context"
	"strings"
	"time"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
)

// AccessStorage Microsoft Access 数据库存储实现
type AccessStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewAccessStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*AccessStorage, error) {
	// 如果调用方没有显式指定方言（还是默认的 Sql92Provider），换成 Access 自己的方言
	// 提供者：Sql92Provider.NowTimestampSql 是 MySQL 专有语法，在 Access 上会直接报语法错误
	if _, ok := options.SqlProvider.(*sql_based_storage.Sql92Provider); ok || options.SqlProvider == nil {
		options.SetSqlProvider(NewAccessSqlProvider())
	}
	baseStorage, err := sql_based_storage.NewSqlBasedStorage(options)
	if err != nil {
		return nil, err
	}
	return &AccessStorage{SqlBasedStorage: baseStorage}, nil
}

// CreateWithVersion Access 唯一索引冲突时返回 ErrLockAlreadyExists
// 注意：未接入真实 Access 驱动验证，错误消息按 Access/ODBC 官方文档做的最佳努力匹配。
func (x *AccessStorage) CreateWithVersion(ctx context.Context, lockId string, version storage.Version, lockInformation *storage.LockInformation) error {
	err := x.SqlBasedStorage.CreateWithVersion(ctx, lockId, version, lockInformation)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "duplicate values in the index") || strings.Contains(msg, "duplicate values") {
			return storage_lock.ErrLockAlreadyExists
		}
	}
	return err
}

const StorageName = "access-storage"

func (x *AccessStorage) GetName() string { return StorageName }

func (x *AccessStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{storage.CapabilityCAS, storage.CapabilityReliableTime, storage.CapabilityWatch}
}

func (x *AccessStorage) Watch(ctx context.Context, lockId string) (<-chan *storage.LockInformation, error) {
	return storage.NewPollingLockWatcher(x, 500*time.Millisecond).Watch(ctx, lockId)
}
