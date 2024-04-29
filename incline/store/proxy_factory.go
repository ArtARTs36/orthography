package store

import "database/sql"

func NewProxyMemoryAndDB(tableName string, db *sql.DB) *Proxy {
	return NewProxy(NewMemory(), NewDB(tableName, db))
}
