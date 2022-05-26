package engines

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cdleo/go-commons/sqlcommons"
)

type mockDBSqlConn struct {
	mock sqlmock.Sqlmock
}

func NewMockDBSqlAdapter() sqlcommons.SQLEngineAdapter {

	return &mockDBSqlConn{}
}

func (s *mockDBSqlConn) Open() (*sql.DB, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	s.mock = mock
	return db, err
}

func (s *mockDBSqlConn) ErrorHandler(err error) error {
	return err
}
