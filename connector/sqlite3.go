package connector

import (
	"database/sql"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/mattn/go-sqlite3"
)

type sqlite3Conn struct {
	url string
}

const sqlite3ProxyName = "sqlite3-proxy"

func NewSqlite3Connector(url string) sqlcommons.SQLConnector {
	return &sqlite3Conn{
		url,
	}
}

func (s *sqlite3Conn) Open(logger logger.Logger, translator sqlcommons.SQLAdapter) (*sql.DB, error) {

	registerProxy(sqlite3ProxyName, logger, translator, &sqlite3.SQLiteDriver{})

	return sql.Open(sqlite3ProxyName, s.url)
}

func (s *sqlite3Conn) GetNextSequenceQuery(sequenceName string) string {
	return sequenceName
}
