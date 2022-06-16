# go-sqldb

[![Go Reference](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/cdleo/go-sqldb) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/cdleo/go-sqldb/master/LICENSE) [![Build Status](https://scrutinizer-ci.com/g/cdleo/go-sqldb/badges/build.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sqldb/build-status/main) [![Code Coverage](https://scrutinizer-ci.com/g/cdleo/go-sqldb/badges/coverage.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sqldb/?branch=main) [![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/cdleo/go-sqldb/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sqldb/?branch=main)

go-sqlDB it's a muti DB Engine wrapper for the GO **database/sql** package. It provides a set of standard error codes, providing abstraction from the implemented DB engine and allowing changing it (almost) just modifying configuration, without the need to modify source code.
Besides that, provides a very limited cross-engine sql translator.

## General
The sqlClient contract resides on the go-commons repository: [github.com/cdleo/go-commons/sqlcommons/sqlClient.go](https://ithub.com/cdleo/go-commons/sqlcommons/sqlClient.go):
```go
type SQLClient interface {
	Open() error
	Close()
	IsOpen() error

	Begin() (SQLTx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (SQLTx, error)

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	Prepare(query string) (SQLStmt, error)
	PrepareContext(ctx context.Context, query string) (SQLStmt, error)

	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
```

**Supported Engines**
Currently, the next set of engines are supported:
- **Oracle**: Using the godror driver [github.com/godror/godror](https://github.com/godror/godror)
- **Postgres**: Using the pq driver [github.com/lib/pq](https://github.com/lib/pq)
- **SQLite3**: Using the go-sqlite3 driver [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)


**Usage**
This example program shows the initialization and the use at basic level:
```go
package sqldb_test

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cdleo/go-sqldb"
	"github.com/cdleo/go-sqldb/engines"
	"github.com/cdleo/go-sqldb/translators"
)

type People struct {
	Id       int    `db:"id"`
	Nombre   string `db:"firstname"`
	Apellido string `db:"lastname"`
}

func Example_sqlConn() {

	adapter := engines.NewSqlite3Adapter(":memory:")
	translator := translators.NewNoopTranslator()
	sqlConn := sqldb.NewSQLDB(adapter, translator)

	if err := sqlConn.Open(); err != nil {
		fmt.Println("Unable to connect to DB")
		os.Exit(1)
	}
	defer sqlConn.Close()

	statement, err := sqlConn.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		fmt.Printf("Unable to prepare statement %v\n", err)
		os.Exit(1)
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Printf("Unable to exec statement %v\n", err)
		os.Exit(1)
	}

	statement, err = sqlConn.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		fmt.Printf("Unable to prepare statement %v\n", err)
		os.Exit(1)
	}
	_, err = statement.Exec("Gene", "Kranz")
	if err != nil {
		fmt.Printf("Unable to exec statement %v\n", err)
		os.Exit(1)
	}

	rows, err := sqlConn.Query("SELECT id, firstname, lastname FROM people")
	if err != nil {
		fmt.Printf("Unable to query data %v\n", err)
		os.Exit(1)
	}

	var p People
	for rows.Next() {
		_ = rows.Scan(&p.Id, &p.Nombre, &p.Apellido)
		fmt.Println(strconv.Itoa(p.Id) + ": " + p.Nombre + " " + p.Apellido)
	}

	// Output:
	// 1: Gene Kranz
}
```

## Sample

You can find a sample of the use of go-sqldb project [HERE](https://github.com/cdleo/go-sqldb/blob/master/sqlDB_example_test.go)

## Contributing

Comments, suggestions and/or recommendations are always welcomed. Please check the [Contributing Guide](CONTRIBUTING.md) to learn how to get started contributing.
