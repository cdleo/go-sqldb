package sqlcommons

import (
	"context"
	"database/sql"
)

type SQLClient interface {
	Open() error
	Close()
	IsOpen() error

	Begin() (SQLTx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (SQLTx, error)

	SQLTransactionable
}

type SQLTransactionable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	Prepare(query string) (SQLStmt, error)
	PrepareContext(ctx context.Context, query string) (SQLStmt, error)

	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type SQLTx interface {
	Commit() error

	SQLTransactionable

	Rollback() error
}

type SQLStmt interface {
	Close() error

	Exec(args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)

	Query(args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error)
	QueryRow(args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row
}
