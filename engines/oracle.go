package engines

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-sqldb"
	"github.com/godror/godror"
	"github.com/godror/godror/dsn"
)

type oracleConn struct {
	connString string
	user       string
	password   string
}

const oracleProxyName = "godror-proxy"

func NewOracleSqlAdapter(host string, port int, user string, password string, database string) sqldb.SQLEngineAdapter {

	return &oracleConn{
		connString: fmt.Sprintf("(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=%s)(PORT=%d))(CONNECT_DATA=(%s)))", host, port, database),
		user:       user,
		password:   password,
	}
}

func NewOracleTNSSqlAdapter(tnsName string, user string, password string) sqldb.SQLEngineAdapter {

	return &oracleConn{
		connString: fmt.Sprintf("connectString=%s", tnsName),
		user:       user,
		password:   password,
	}
}

func (s *oracleConn) Open(logger logger.Logger, translator sqldb.SQLSyntaxTranslator) (*sql.DB, error) {

	registerProxy(oracleProxyName, logger, translator, godror.NewConnector(dsn.ConnectionParams{}).Driver())

	return sql.Open(oracleProxyName, godror.ConnectionParams{
		CommonParams: godror.CommonParams{
			ConnectString: s.connString,
			Username:      s.user,
			Password:      godror.NewPassword(s.password),
			Timezone:      time.Local,
		},
		StandaloneConnection: true,
	}.StringWithPassword())
}
