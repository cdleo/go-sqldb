package engines

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/cdleo/go-commons/logger"
	proxy "github.com/cdleo/go-sql-proxy"
	"github.com/cdleo/go-sqldb"
)

func registerProxy(name string, logger logger.Logger, translator sqldb.SQLSyntaxTranslator, sqlDriver driver.Driver) {

	drivers := sql.Drivers()
	for _, item := range drivers {
		if item == name {
			return
		}
	}

	sql.Register(name, proxy.NewProxyContext(sqlDriver, &proxy.HooksContext{
		Open: func(_ context.Context, _ interface{}, conn *proxy.Conn) error {
			logger.Qry("Open conn")
			return nil
		},
		Close: func(_ context.Context, _ interface{}, conn *proxy.Conn) error {
			logger.Qry("Close conn")
			return nil
		},

		PreExec: func(_ context.Context, stmt *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return time.Now(), nil
		},
		PostExec: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Result, _ error) error {
			logger.Tracef("Exec: %s; args = %v (%s)", prettyQuery(stmt.QueryString), args, time.Since(ctx.(time.Time)))
			return nil
		},

		PreQuery: func(_ context.Context, stmt *proxy.Stmt, _ []driver.NamedValue) (interface{}, error) {
			stmt.QueryString = translator.Translate(stmt.QueryString)
			return time.Now(), nil
		},
		PostQuery: func(_ context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, _ driver.Rows, _ error) error {
			logger.Tracef("Query: %s; args = %v (%s)", prettyQuery(stmt.QueryString), args, time.Since(ctx.(time.Time)))
			return nil
		},

		Begin: func(_ context.Context, _ interface{}, conn *proxy.Conn) error {
			logger.Qry("Begin")
			return nil
		},
		Commit: func(_ context.Context, _ interface{}, tx *proxy.Tx) error {
			logger.Qry("Commit")
			return nil
		},
		Rollback: func(_ context.Context, _ interface{}, tx *proxy.Tx) error {
			logger.Qry("Rollback")
			return nil
		},

		OnError: func(ctx interface{}, err error) error {
			return translator.ErrorHandler(err)
		},
	}))
}

func prettyQuery(query string) string {
	return strings.ReplaceAll(strings.ReplaceAll(query, "\t", ""), "\n", "")
}
