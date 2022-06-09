package sqldb

import (
	"database/sql"
	"database/sql/driver"
)

//Implemented Engines
const (
	Oracle_Engine     = "Oracle"
	PostgreSQL_Engine = "PostgreSQL"
	SQLite3_Engine    = "SQLite3"
	MockDB_Engine     = "MockDB"
)

type MockSQLEngineAdapter interface {
	SQLEngineAdapter

	PatchBegin(err error)
	PatchCommit(err error)
	PatchRollback(err error)

	PatchExec(query string, err error, args ...driver.Value)
	PatchQuery(query string, columns []string, values []driver.Value, err error, args ...driver.Value)
	PatchQueryRow(query string, result map[string]string, err error)
}

//go:generate mockgen -package enginesMocks -destination engines/mocks/sqlEngineAdapter.go . SQLEngineAdapter
type SQLEngineAdapter interface {
	Open() (*sql.DB, error)
	ErrorHandler(err error) error
}

//go:generate mockgen -package translatorsMocks -destination translators/mocks/sqlSyntaxTranslator.go . SQLSyntaxTranslator
type SQLSyntaxTranslator interface {
	Translate(query string) string
}
