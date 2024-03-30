package sqlproxy

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
