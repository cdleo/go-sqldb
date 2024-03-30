package connector

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/godror/godror"
	"github.com/godror/godror/dsn"
)

type oracleConn struct {
	connString string
	user       string
	password   string
}

const oracleProxyName = "godror-proxy"

func NewOracleSqlConnector(host string, port int, user string, password string, database string) sqlcommons.SQLConnector {

	return &oracleConn{
		connString: fmt.Sprintf("(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=%s)(PORT=%d))(CONNECT_DATA=(%s)))", host, port, database),
		user:       user,
		password:   password,
	}
}

func NewOracleTNSSqlConnector(tnsName string, user string, password string) sqlcommons.SQLConnector {

	return &oracleConn{
		connString: fmt.Sprintf("connectString=%s", tnsName),
		user:       user,
		password:   password,
	}
}

func (s *oracleConn) Open(logger logger.Logger, translator sqlcommons.SQLAdapter) (*sql.DB, error) {

	registerProxy(oracleProxyName, logger, translator, godror.NewConnector(dsn.ConnectionParams{}).Driver())

	var connParams godror.ConnectionParams
	connParams.ConnectString = s.connString
	connParams.Username = s.user
	connParams.Password = godror.NewPassword(s.password)
	connParams.Timezone = time.Local
	connParams.StandaloneConnection = true

	return sql.Open(oracleProxyName, connParams.StringWithPassword())
}

func (s *oracleConn) GetNextSequenceQuery(sequenceName string) string {
	return fmt.Sprintf("SELECT %s.NEXTVAL FROM DUAL", sequenceName)
}
