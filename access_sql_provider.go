package access_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
)

// AccessSqlProvider Microsoft Access（Jet/ACE 引擎）的 SQL 方言提供者
// Access 不支持 MySQL 的 UNIX_TIMESTAMP(NOW())。
type AccessSqlProvider struct {
	*sql_based_storage.Sql92Provider
}

func NewAccessSqlProvider() *AccessSqlProvider {
	return &AccessSqlProvider{
		Sql92Provider: sql_based_storage.NewSql92Provider(),
	}
}

// NowTimestampSql Access 获取当前时间的 Unix 时间戳
// 修复：原来直接继承 Sql92Provider 的 "SELECT UNIX_TIMESTAMP(NOW())"，Access 不支持这两个 MySQL
// 专有函数，执行即报语法错误，导致 GetTime 完全不可用。用 Access 的 DATEDIFF('s', ...) 函数
// 换算 1970-01-01 到 NOW() 的秒数差作为 Unix 时间戳。
// 注意：未接入真实 Access 驱动（ODBC）验证，是按 Access SQL 官方文档做的最佳努力实现。
func (x *AccessSqlProvider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return `SELECT DATEDIFF('s', #1970-01-01 00:00:00#, NOW())`, nil
}
