package sqldb

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
