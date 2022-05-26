package engines

import (
	"database/sql"
	"fmt"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/lib/pq"
)

type pgSqlConn struct {
	url      string
	user     string
	password string
	database string
}

const postgresql_DriverName = "postgres"

func NewPostgreSqlAdapter(host string, port int, user string, password string, database string) sqlcommons.SQLEngineAdapter {

	return &pgSqlConn{
		url:      fmt.Sprintf("%s:%d", host, port),
		user:     user,
		password: password,
		database: database,
	}
}

func (s *pgSqlConn) Open() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", s.user, s.password, s.url, s.database)
	return sql.Open(postgresql_DriverName, dataSourceName)
}

func (s *pgSqlConn) ErrorHandler(err error) error {
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
