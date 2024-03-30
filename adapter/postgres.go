package adapter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/lib/pq"
)

type postgresAdapter struct {
	paramRegExp     *regexp.Regexp
	sourceSQLSintax string
}

func NewPostgresAdapter(sourceSQLSintax string) sqlcommons.SQLAdapter {
	return &postgresAdapter{
		regexp.MustCompile(":[1-9]"),
		sourceSQLSintax,
	}
}

func (s *postgresAdapter) Translate(query string) string {

	if s.sourceSQLSintax == "Oracle" {
		return s.paramRegExp.ReplaceAllStringFunc(query, func(m string) string {
			return strings.Replace(m, ":", "$", 1)
		})
	} else {
		return query
	}
}

func (s *postgresAdapter) ErrorHandler(err error) error {
	if err == nil {
		return nil
	}

	if pqError, ok := err.(*pq.Error); ok {
		switch pqError.Code {
		case "23505":
			return sqlcommons.UniqueConstraintViolation
		case "23503":
			return sqlcommons.IntegrityConstraintViolation
		case "22001":
			return sqlcommons.ValueTooLargeForColumn
		case "22003":
			return sqlcommons.ValueLargerThanPrecision
		case "23502":
			return sqlcommons.CannotSetNullColumn
		case "22P02":
			return sqlcommons.InvalidNumericValue
		case "21000":
			return sqlcommons.SubqueryReturnsMoreThanOneRow
		default:
			return fmt.Errorf("Unhandled PostgreSQL error. Code:[%s] Desc:[%s]", pqError.Code, pqError.Message)
		}
	} else {
		return err
	}
}
