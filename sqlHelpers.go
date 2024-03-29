package sqldb

import (
	"database/sql"
	"database/sql/driver"

	"github.com/cdleo/go-commons/logger"
)

type DBEngine string

const (
	Oracle     DBEngine = "Oracle"
	PostgreSQL DBEngine = "PostgreSQL"
	SQLite3    DBEngine = "SQLite3"
	MockDB     DBEngine = "MockDB"
)

var DBEngines = []DBEngine{
	Oracle,
	PostgreSQL,
	SQLite3,
	MockDB,
}

type SQLSintaxTranslator string

const (
	None         SQLSintaxTranslator = "None"
	ToOracle     SQLSintaxTranslator = "ToOracle"
	ToPostgreSQL SQLSintaxTranslator = "ToPostgreSQL"
	ToSQLite3    SQLSintaxTranslator = "ToSQLite3"
)

var SQLSintaxTranslators = []SQLSintaxTranslator{
	None,
	ToOracle,
	ToPostgreSQL,
	ToSQLite3,
}

type MockSQLEngineAdapter interface {
	SQLEngineAdapter

	PatchBegin(err error)
	PatchCommit(err error)
	PatchRollback(err error)

	PatchExec(query string, err error, args ...driver.Value)
	PatchQuery(query string, columns []string, values []driver.Value, err error, args ...driver.Value)
	PatchQueryRow(query string, result map[string]string, err error)
}

type SQLEngineAdapter interface {
	Open(logger logger.Logger, translator SQLSyntaxTranslator) (*sql.DB, error)
}

//go:generate mockgen -package translatorsMocks -destination translators/mocks/sqlSyntaxTranslator.go . SQLSyntaxTranslator
type SQLSyntaxTranslator interface {
	Translate(query string) string
	ErrorHandler(err error) error
}
