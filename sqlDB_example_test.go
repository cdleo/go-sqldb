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
