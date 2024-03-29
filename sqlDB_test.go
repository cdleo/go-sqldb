package sqldb_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/cdleo/go-sqldb"
	"github.com/cdleo/go-sqldb/engines"

	"github.com/cdleo/go-sqldb/translator"
	"github.com/stretchr/testify/require"
)

type Customers struct {
	Id         int           `db:"id"`
	Name       string        `db:"name"`
	Updatetime time.Time     `db:"updatetime"`
	Age        sql.NullInt64 `db:"age"`
	Group      int           `db:"cust_group"`
	Dummy      string        `db:"not_existing_field"`
}

func Test_sqlConn_InitErr(t *testing.T) {
	// Setup
	adapter := engines.NewMockSQLAdapter(false)
	translator := translator.NewNoopTranslator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)
	_, err := sqlConn.Open()
	require.Error(t, err)
}

func Test_sqlConn_InitOK(t *testing.T) {
	// Setup
	adapter := engines.NewMockSQLAdapter(true)
	translator := translator.NewNoopTranslator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)
	_, err := sqlConn.Open()
	require.NoError(t, err)
}

func Test_sqlConn_CreateTables(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	// Exec
	require.NoError(t, createTablesHelper(sqlDB))

	sqlConn.Close()
}

func Test_sqlConn_DropTables(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)
	require.NoError(t, createTablesHelper(sqlDB))

	// Exec
	require.NoError(t, dropTablesHelper(sqlDB))
}

func Test_sqlConn_StoreData(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))

	// Exec
	require.NoError(t, insertDataHelper(sqlDB))

	sqlConn.Close()
}

func Test_sqlConn_ReturnData(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))
	require.NoError(t, insertDataHelper(sqlDB))

	// Exec
	rows, err2 := sqlDB.Query("SELECT name FROM customers")
	defer rows.Close()

	require.NoError(t, err2)
	require.True(t, rows.Next())

	sqlConn.Close()
}

func Test_sqlConn_CanThrowInvalidTableError(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))
	require.NoError(t, insertDataHelper(sqlDB))

	// Exec
	_, err2 := sqlDB.Query("SELECT name FROM customerxs")

	require.Error(t, err2)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowCannotInsertNullError(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))

	// Exec
	_, err2 := sqlDB.Exec("INSERT INTO customers (name, updatetime) VALUES (:1,:2)", nil, time.Now())

	require.ErrorIs(t, err2, sqlcommons.CannotSetNullColumn)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowCannotUpdateNullError(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()
	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))
	require.NoError(t, insertDataHelper(sqlDB))

	// Exec
	_, err2 := sqlDB.Exec("UPDATE customers c SET name = :1 WHERE c.name = :2", nil, "Pablo")

	require.Error(t, err2, sqlcommons.CannotSetNullColumn)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowUniqueConstraintViolationError(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:?_foreign_keys=on")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))
	require.NoError(t, insertDataHelper(sqlDB))

	// Exec
	_, err2 := sqlDB.Exec("INSERT INTO customers (name, updatetime, age, cust_group)VALUES(:1, :2, :3, :4)", "Juan", time.Now(), nil, 1)

	require.ErrorIs(t, err2, sqlcommons.UniqueConstraintViolation)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowForeignKeyConstraintViolationError(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:?_foreign_keys=on")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))
	require.NoError(t, insertDataHelper(sqlDB))

	// Exec
	_, err2 := sqlDB.Exec("UPDATE customers SET cust_group = :1 WHERE name = :2", 2, "Pablo")

	require.ErrorIs(t, err2, sqlcommons.IntegrityConstraintViolation)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowIntegrityConstraintViolationError(t *testing.T) {
	// Setup
	adapter := engines.NewSqlite3Adapter(":memory:?_foreign_keys=on")
	translator := translator.NewSQLite3Translator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	require.NoError(t, createTablesHelper(sqlDB))
	require.NoError(t, insertDataHelper(sqlDB))

	// Exec
	_, err2 := sqlDB.Exec("DELETE from customers_groups WHERE id = :1", 1)

	require.ErrorIs(t, err2, sqlcommons.IntegrityConstraintViolation)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowValueTooLargeError(t *testing.T) {
	// Setup
	adapter := engines.NewMockSQLAdapter(true)
	translator := translator.NewNoopTranslator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	// Exec
	query := fmt.Sprintf("INSERT INTO customers (name, cust_group) VALUES (%s,%d)", "'verylongname'", 1)
	adapter.PatchExec(query, sqlcommons.ValueTooLargeForColumn)

	_, err2 := sqlDB.Exec(query)

	require.ErrorIs(t, err2, sqlcommons.ValueTooLargeForColumn)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowSubqueryReturnsMoreThanOneRowError(t *testing.T) {
	// Setup
	adapter := engines.NewMockSQLAdapter(true)
	translator := translator.NewNoopTranslator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	query := "SELECT name FROM customers WHERE id = (SELECT id FROM customers)"
	adapter.PatchQuery(query, nil, nil, sqlcommons.SubqueryReturnsMoreThanOneRow)

	// Exec
	_, err2 := sqlDB.Query(query)

	require.ErrorIs(t, err2, sqlcommons.SubqueryReturnsMoreThanOneRow)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowInvalidNumericValueError(t *testing.T) {
	// Setup
	adapter := engines.NewMockSQLAdapter(true)
	translator := translator.NewNoopTranslator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	query := "UPDATE customers SET age = :1 WHERE name = :2"
	adapter.PatchExec(query, sqlcommons.InvalidNumericValue, "twelve", "Pablo")

	// Exec
	_, err2 := sqlDB.Exec(query, "twelve", "Pablo")

	require.ErrorIs(t, err2, sqlcommons.InvalidNumericValue)

	sqlConn.Close()
}

func Test_sqlConn_CanThrowValueLargerThanPrecisionError(t *testing.T) {
	// Setup
	adapter := engines.NewMockSQLAdapter(true)
	translator := translator.NewNoopTranslator()
	logger := sqldb.NewNoLogLogger()

	sqlConn := sqldb.NewSQLDB(adapter, translator, logger)

	sqlDB, err := sqlConn.Open()
	require.NoError(t, err)

	query := "UPDATE customers SET age = :1 WHERE name = :2"
	adapter.PatchExec(query, sqlcommons.ValueLargerThanPrecision, 949.0044, "Pablo")

	// Exec
	_, err2 := sqlDB.Exec(query, 949.0044, "Pablo")

	require.ErrorIs(t, err2, sqlcommons.ValueLargerThanPrecision)

	sqlConn.Close()
}

func createTablesHelper(sqlClient *sql.DB) error {

	if _, err := sqlClient.Exec(`CREATE TABLE IF NOT EXISTS customers_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		groupname TEXT NOT NULL)`); err != nil {
		return err
	}
	if _, err := sqlClient.Exec(`CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name CHAR(10) NOT NULL,
		updatetime TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
		age INT NULL,
		cust_group INT NOT NULL,
		FOREIGN KEY (cust_group) REFERENCES customers_groups (id) ON DELETE RESTRICT
		CONSTRAINT customers_un UNIQUE (name))`); err != nil {
		return err
	}
	return nil
}

func insertDataHelper(sqlClient *sql.DB) error {

	if _, err := sqlClient.Exec(`INSERT INTO customers_groups (groupname) VALUES('General');`); err != nil {
		return err
	}

	if statement, err := sqlClient.Prepare("INSERT INTO customers (name, updatetime, age, cust_group)VALUES(:1, :2, :3, :4)"); err != nil {
		return err
	} else {
		if _, err := statement.Exec("Juan", time.Now(), nil, 1); err != nil {
			return err
		}
		if _, err := statement.Exec("Pedro", time.Now(), nil, 1); err != nil {
			return err
		}
		if _, err := statement.Exec("Pablo", time.Now(), 99, 1); err != nil {
			return err
		}
	}

	return nil
}

func dropTablesHelper(sqlClient *sql.DB) error {

	if _, err := sqlClient.Exec(`DROP TABLE IF EXISTS customers_groups`); err != nil {
		return err
	}
	if _, err := sqlClient.Exec(`DROP TABLE IF EXISTS customers`); err != nil {
		return err
	}
	return nil
}
