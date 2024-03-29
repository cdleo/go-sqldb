package translator

import (
	"fmt"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/cdleo/go-sqldb"
	"github.com/mattn/go-sqlite3"
)

type sqlite3Translator struct{}

func NewSQLite3Translator() sqldb.SQLSyntaxTranslator {
	return &sqlite3Translator{}
}

func (t *sqlite3Translator) Translate(query string) string {
	return query
}

func (s *sqlite3Translator) ErrorHandler(err error) error {

	if sqliteError, ok := err.(sqlite3.Error); ok {

		if sqliteError.Code == 18 { //SQLITE_TOOBIG
			return sqlcommons.ValueTooLargeForColumn

		} else if sqliteError.Code == 19 { //SQLITE_CONSTRAINT
			if sqliteError.ExtendedCode == 787 || /*SQLITE_CONSTRAINT_FOREIGNKEY*/
				sqliteError.ExtendedCode == 1555 || /*SQLITE_CONSTRAINT_PRIMARYKEY*/
				sqliteError.ExtendedCode == 1811 { /*SQLITE_CONSTRAINT_TRIGGER*/
				return sqlcommons.IntegrityConstraintViolation

			} else if sqliteError.ExtendedCode == 1299 { //SQLITE_CONSTRAINT_NOTNULL
				return sqlcommons.CannotSetNullColumn

			} else if sqliteError.ExtendedCode == 2067 { //SQLITE_CONSTRAINT_UNIQUE
				return sqlcommons.UniqueConstraintViolation

			}
		} else if sqliteError.Code == 25 { //SQLITE_RANGE
			return sqlcommons.InvalidNumericValue
		}

		return fmt.Errorf("Unhandled SQLite3 error. Code:[%s] Extended:[%s] Desc:[%s]", sqliteError.Code, sqliteError.ExtendedCode, sqliteError.Error())

	} else {
		return err
	}
}
