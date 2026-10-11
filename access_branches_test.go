package access_storage

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

type stubProvider struct{ *sql_based_storage.Sql92Provider }

func (stubProvider) NowTimestampSql(context.Context, string) (string, []any) {
	return "SELECT stub_access_time", nil
}

func TestAccessSqlProvider(t *testing.T) {
	sqlText, args := NewAccessSqlProvider().NowTimestampSql(context.Background(), "locks")
	assert.Contains(t, sqlText, "DATEDIFF")
	assert.Nil(t, args)
}

func TestNewAccessStorage_NilOptions(t *testing.T) {
	assert.Panics(t, func() { _, _ = NewAccessStorage(context.Background(), nil) })
}

func TestNewAccessStorage_InvalidOptions(t *testing.T) {
	_, err := NewAccessStorage(context.Background(), &sql_based_storage.SqlBasedStorageOptions{})
	assert.Error(t, err)
}

func TestNewAccessStorage_CustomProvider(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	opts := sql_based_storage.NewSqlBasedStorageOptions().
		SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db)).
		SetSqlProvider(stubProvider{Sql92Provider: sql_based_storage.NewSql92Provider()})
	s, err := NewAccessStorage(context.Background(), opts)
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT stub_access_time").WillReturnRows(sqlmock.NewRows([]string{"ts"}).AddRow(uint64(1)))
	_, err = s.GetTime(context.Background())
	assert.NoError(t, err)
}

func TestAccessStorage_CreateWithVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	opts := sql_based_storage.NewSqlBasedStorageOptions().
		SetConnectionManager(storage.NewFixedSqlDBConnectionManager(db))
	s, err := NewAccessStorage(context.Background(), opts)
	assert.NoError(t, err)

	info := &storage.LockInformation{OwnerId: "owner"}
	mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("duplicate values in the index"))
	err = s.CreateWithVersion(context.Background(), "lock", 1, info)
	assert.ErrorIs(t, err, storage_lock.ErrLockAlreadyExists)

	mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("duplicate values"))
	err = s.CreateWithVersion(context.Background(), "lock", 1, info)
	assert.ErrorIs(t, err, storage_lock.ErrLockAlreadyExists)

	mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("other"))
	err = s.CreateWithVersion(context.Background(), "lock", 1, info)
	assert.EqualError(t, err, "other")

	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	assert.NoError(t, s.CreateWithVersion(context.Background(), "lock", 1, info))
}

func TestAccessStorage_Watch_EmptyLockID(t *testing.T) {
	s, err := NewAccessStorage(context.Background(), sql_based_storage.NewSqlBasedStorageOptions().
		SetConnectionManager(storage.NewFixedSqlDBConnectionManager(&sql.DB{})))
	assert.NoError(t, err)
	_, err = s.Watch(context.Background(), "")
	assert.Error(t, err)
}
