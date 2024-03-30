package connector

import (
	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-commons/sqlcommons"
)

type mockDBSqlConn struct {
	initOk bool
	mock   sqlmock.Sqlmock
}

func NewMockSQLConnector(initOk bool) sqlcommons.MockSQLConnector {

	return &mockDBSqlConn{
		initOk,
		nil,
	}
}

func (s *mockDBSqlConn) Open(logger logger.Logger, translator sqlcommons.SQLAdapter) (*sql.DB, error) {

	if s.initOk {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		s.mock = mock
		return db, err
	} else {
		return nil, sqlcommons.ConnectionFailed
	}
}

func (s *mockDBSqlConn) GetNextSequenceQuery(sequenceName string) string {
	return sequenceName
}

func (s *mockDBSqlConn) PatchBegin(err error) {
	expectBegin := s.mock.ExpectBegin()
	if err != nil {
		expectBegin.WillReturnError(err)
	}
}
func (s *mockDBSqlConn) PatchCommit(err error) {
	expectCommit := s.mock.ExpectCommit()
	if err != nil {
		expectCommit.WillReturnError(err)
	}
}
func (s *mockDBSqlConn) PatchRollback(err error) {
	expectRollback := s.mock.ExpectRollback()
	if err != nil {
		expectRollback.WillReturnError(err)
	}
}

func (s *mockDBSqlConn) PatchExec(query string, err error, args ...driver.Value) {
	expectQuery := s.mock.ExpectExec(query)
	if len(args) > 0 {
		expectQuery.WithArgs(args...)
	}
	if err != nil {
		expectQuery.WillReturnError(err)
	} else {

		result := sqlmock.NewResult(0, 0)
		expectQuery.WillReturnResult(result)
	}
}
func (s *mockDBSqlConn) PatchQuery(query string, columns []string, values []driver.Value, err error, args ...driver.Value) {
	expectQuery := s.mock.ExpectQuery(query)
	if len(args) > 0 {
		expectQuery.WithArgs(args...)
	}
	if err != nil {
		expectQuery.WillReturnError(err)
	} else {
		rows := sqlmock.NewRows(columns).
			AddRow(values...)

		expectQuery.WillReturnRows(rows)
	}
}

func (s *mockDBSqlConn) PatchQueryRow(query string, result map[string]string, err error) {
	var keys []string
	var values []interface{}
	for key, value := range result {
		keys = append(keys, key)
		values = append(values, value)
	}
	rows := sqlmock.NewRows(keys).
		AddRow(values)

	s.mock.ExpectQuery(query).WillReturnRows(rows).WillReturnError(err)
}
