package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/cdleo/go-commons/sqlcommons"
)

type sqlDB struct {
	engineAdapter sqlcommons.SQLEngineAdapter
	translator    sqlcommons.SQLSyntaxTranslator
	db            *sql.DB
}

type sqlTx struct {
	adapter    sqlcommons.SQLEngineAdapter
	translator sqlcommons.SQLSyntaxTranslator
	tx         *sql.Tx
}

type sqlStmt struct {
	adapter sqlcommons.SQLEngineAdapter
	stmt    *sql.Stmt
}

func NewSQLDB(adapter sqlcommons.SQLEngineAdapter, translator sqlcommons.SQLSyntaxTranslator) sqlcommons.SQLClient {
	return &sqlDB{
		adapter,
		translator,
		nil,
	}
}

func (s *sqlDB) Open() error {
	if db, err := s.engineAdapter.Open(); err != nil {
		return s.engineAdapter.ErrorHandler(err)
	} else {
		s.db = db
	}
	return nil
}

func (s *sqlDB) Close() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
}

func (s *sqlDB) IsOpen() error {
	if s.db == nil {
		return sqlcommons.DBNotInitialized
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if stdErr := s.db.PingContext(ctx); stdErr != nil {
		s.Close()
		if err := s.Open(); err != nil {
			return s.engineAdapter.ErrorHandler(err)
		} else {
			ctxReconnect, cancelReconnect := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancelReconnect()
			return s.db.PingContext(ctxReconnect)
		}
	}
	return nil
}

func (s *sqlDB) Begin() (sqlcommons.SQLTx, error) {
	if err := s.IsOpen(); err != nil {
		return nil, err
	}
	if tx, err := s.db.Begin(); err != nil {
		return nil, err
	} else {
		return newSQLTx(tx, s.translator, s.engineAdapter), nil
	}
}
func (s *sqlDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqlcommons.SQLTx, error) {
	if err := s.IsOpen(); err != nil {
		return nil, err
	}
	if tx, err := s.db.BeginTx(ctx, opts); err != nil {
		return nil, err
	} else {
		return newSQLTx(tx, s.translator, s.engineAdapter), nil
	}
}

func (s *sqlDB) Prepare(query string) (sqlcommons.SQLStmt, error) {
	if err := s.IsOpen(); err != nil {
		return nil, err
	}
	if stmt, err := s.db.Prepare(s.translator.Translate(query)); err != nil {
		return nil, s.engineAdapter.ErrorHandler(err)
	} else {
		return newSQLStmt(stmt, s.engineAdapter), nil
	}
}
func (s *sqlDB) PrepareContext(ctx context.Context, query string) (sqlcommons.SQLStmt, error) {
	if err := s.IsOpen(); err != nil {
		return nil, err
	}
	if stmt, err := s.db.PrepareContext(ctx, s.translator.Translate(query)); err != nil {
		return nil, s.engineAdapter.ErrorHandler(err)
	} else {
		return newSQLStmt(stmt, s.engineAdapter), nil
	}
}

func (s *sqlDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.Exec(s.translator.Translate(query), args...)
	return result, s.engineAdapter.ErrorHandler(err)
}
func (s *sqlDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.ExecContext(ctx, s.translator.Translate(query), args...)
	return result, s.engineAdapter.ErrorHandler(err)
}

func (s *sqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	result, err := s.db.Query(s.translator.Translate(query), args...)
	return result, s.engineAdapter.ErrorHandler(err)
}
func (s *sqlDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	result, err := s.db.QueryContext(ctx, s.translator.Translate(query), args...)
	return result, s.engineAdapter.ErrorHandler(err)
}
func (s *sqlDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(s.translator.Translate(query), args...)
}
func (s *sqlDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, s.translator.Translate(query), args...)
}

/* ---------------------------------------------------------------------------------------------------------------------- */
func newSQLTx(tx *sql.Tx, translator sqlcommons.SQLSyntaxTranslator, adapter sqlcommons.SQLEngineAdapter) sqlcommons.SQLTx {
	return &sqlTx{
		adapter,
		translator,
		tx,
	}
}

func (s *sqlTx) Commit() error {
	return s.adapter.ErrorHandler(s.tx.Commit())
}

func (s *sqlTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.tx.Exec(s.translator.Translate(query), args...)
	return result, s.adapter.ErrorHandler(err)
}
func (s *sqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := s.tx.ExecContext(ctx, s.translator.Translate(query), args...)
	return result, s.adapter.ErrorHandler(err)
}

func (s *sqlTx) Prepare(query string) (sqlcommons.SQLStmt, error) {
	if stmt, err := s.tx.Prepare(s.translator.Translate(query)); err != nil {
		return nil, s.adapter.ErrorHandler(err)
	} else {
		return newSQLStmt(stmt, s.adapter), nil
	}
}
func (s *sqlTx) PrepareContext(ctx context.Context, query string) (sqlcommons.SQLStmt, error) {
	if stmt, err := s.tx.PrepareContext(ctx, s.translator.Translate(query)); err != nil {
		return nil, s.adapter.ErrorHandler(err)
	} else {
		return newSQLStmt(stmt, s.adapter), nil
	}
}

func (s *sqlTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	result, err := s.tx.Query(s.translator.Translate(query), args...)
	return result, s.adapter.ErrorHandler(err)
}
func (s *sqlTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	result, err := s.tx.QueryContext(ctx, s.translator.Translate(query), args...)
	return result, s.adapter.ErrorHandler(err)
}
func (s *sqlTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.tx.QueryRow(s.translator.Translate(query), args...)
}
func (s *sqlTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.tx.QueryRowContext(ctx, s.translator.Translate(query), args...)
}

func (s *sqlTx) Rollback() error {
	return s.adapter.ErrorHandler(s.tx.Rollback())
}

/* ---------------------------------------------------------------------------------------------------------------------- */
func newSQLStmt(stmt *sql.Stmt, adapter sqlcommons.SQLEngineAdapter) sqlcommons.SQLStmt {
	return &sqlStmt{
		adapter,
		stmt,
	}
}

func (s *sqlStmt) Close() error {
	return s.adapter.ErrorHandler(s.stmt.Close())
}

func (s *sqlStmt) Exec(args ...interface{}) (sql.Result, error) {
	result, err := s.stmt.Exec(args...)
	return result, s.adapter.ErrorHandler(err)
}
func (s *sqlStmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	result, err := s.stmt.ExecContext(ctx, args...)
	return result, s.adapter.ErrorHandler(err)
}

func (s *sqlStmt) Query(args ...interface{}) (*sql.Rows, error) {
	result, err := s.stmt.Query(args...)
	return result, s.adapter.ErrorHandler(err)
}
func (s *sqlStmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	result, err := s.stmt.QueryContext(ctx, args...)
	return result, s.adapter.ErrorHandler(err)
}
func (s *sqlStmt) QueryRow(args ...interface{}) *sql.Row {
	return s.stmt.QueryRow(args...)
}
func (s *sqlStmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	return s.stmt.QueryRowContext(ctx, args...)
}
