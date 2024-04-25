# go-sqldb

[![Go Reference](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/cdleo/go-sqldb) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/cdleo/go-sqldb/master/LICENSE) [![Build Status](https://scrutinizer-ci.com/g/cdleo/go-sqldb/badges/build.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sqldb/build-status/main) [![Code Coverage](https://scrutinizer-ci.com/g/cdleo/go-sqldb/badges/coverage.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sqldb/?branch=main) [![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/cdleo/go-sqldb/badges/quality-score.png?b=main)](https://scrutinizer-ci.com/g/cdleo/go-sqldb/?branch=main)

go-sqlDB it's a muti DB Engine proxy for the GO **database/sql** package. It provides a set of standard error codes, providing abstraction from the implemented DB engine and allowing to switch it without modify the source code, just modifying configuration (almost).
Besides that, provides a VERY LIMITED cross-engine sql syntax translator.

## General
The sqlProxy it's created by the sqlProxyBuilder. The Open() function returns an standard *sql.DB implementation, 
but using a proxy connector to the selected engine.

**Supported Engines**
Currently, the next set of engines are supported:
- **Oracle**: Using the godror driver [github.com/godror/godror](https://github.com/godror/godror)
- **Postgres**: Using the pq driver [github.com/lib/pq](https://github.com/lib/pq)
- **SQLite3**: Using the go-sqlite3 driver [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)


**Usage**
This example program shows initialization and usage at basic level:
```go
package sqldb

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-sqldb/adapter"
	"github.com/cdleo/go-sqldb/connector"
)

type People struct {
	Id       int    `db:"id"`
	Nombre   string `db:"firstname"`
	Apellido string `db:"lastname"`
}

func Example_sqlDBBuilder() {

	sqlProxy := NewSQLProxyBuilder(connector.NewSqlite3Connector(":memory:")).
		WithAdapter(adapter.NewSQLite3Adapter()).
		WithLogger(logger.NewNoLogLogger()).
		Build()

	var sqlDB *sql.DB
	var err error

	if sqlDB, err = sqlProxy.Open(); err != nil {
		fmt.Println("Unable to connect to DB")
		os.Exit(1)
	}
	defer sqlProxy.Close()

	statement, err := sqlDB.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		fmt.Printf("Unable to prepare statement %v\n", err)
		os.Exit(1)
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Printf("Unable to exec statement %v\n", err)
		os.Exit(1)
	}

	statement, err = sqlDB.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		fmt.Printf("Unable to prepare statement %v\n", err)
		os.Exit(1)
	}
	_, err = statement.Exec("Gene", "Kranz")
	if err != nil {
		fmt.Printf("Unable to exec statement %v\n", err)
		os.Exit(1)
	}

	rows, err := sqlDB.Query("SELECT id, firstname, lastname FROM people")
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
