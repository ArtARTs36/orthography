package store

import "database/sql"

func NewProxyMemoryAndDB(tableName string, db *sql.DB) Store {
	return NewProxy(NewMemory(), NewDB(tableName, db))
}
