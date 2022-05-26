package sqldb

import "database/sql"

//Implemented Engines
const (
	Oracle_Engine     = "Oracle"
	PostgreSQL_Engine = "PostgreSQL"
	SQLite3_Engine    = "SQLite3"
	MockDB_Engine     = "MockDB"
)

//go:generate mockgen -package enginesMocks -destination engines/mocks/sqlEngineAdapter.go . SQLEngineAdapter
type SQLEngineAdapter interface {
	Open() (*sql.DB, error)
	ErrorHandler(err error) error
}

//go:generate mockgen -package translatorsMocks -destination translators/mocks/sqlSyntaxTranslator.go . SQLSyntaxTranslator
type SQLSyntaxTranslator interface {
	Translate(query string) string
}
