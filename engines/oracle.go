package engines

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/godror/godror"
)

type oracleConn struct {
	connString string
	user       string
	password   string
}

const oracle_DriverName = "godror"

func NewOracleSqlConn(host string, port int, user string, password string, database string) sqlcommons.SQLEngineAdapter {

	return &oracleConn{
		connString: fmt.Sprintf("(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=%s)(PORT=%d))(CONNECT_DATA=(%s)))", host, port, database),
		user:       user,
		password:   password,
	}
}

func NewOracleTNSSqlConn(tnsName string, user string, password string) sqlcommons.SQLEngineAdapter {

	return &oracleConn{
		connString: fmt.Sprintf("connectString=%s", tnsName),
		user:       user,
		password:   password,
	}
}

func (s *oracleConn) Open() (*sql.DB, error) {

	return sql.Open(oracle_DriverName, godror.ConnectionParams{
		CommonParams: godror.CommonParams{
			ConnectString: s.connString,
			Username:      s.user,
			Password:      godror.NewPassword(s.password),
			Timezone:      time.Local,
		},
		StandaloneConnection: true,
	}.StringWithPassword())
}

func (s *oracleConn) ErrorHandler(err error) error {
	if err == nil {
		return nil
	}

	if oraError, ok := godror.AsOraErr(err); ok {
		switch oraError.Code() {
		case 1: //ORA-00001"
			return sqlcommons.UniqueConstraintViolation
		case 2291, 2292: //ORA-02291 (PKNotFound) AND ORA-02292 (ChildFound)
			return sqlcommons.IntegrityConstraintViolation
		case 12899: //ORA-12899
			return sqlcommons.ValueTooLargeForColumn
		case 1438: //ORA-01438
			return sqlcommons.ValueLargerThanPrecision
		case 1400, 1407: //ORA-01400 (cannot insert) AND ORA-01407 (cannot change value to)
			return sqlcommons.CannotSetNullColumn
		case 1722: //ORA-01722
			return sqlcommons.InvalidNumericValue
		case 1427: //ORA-01427
			return sqlcommons.SubqueryReturnsMoreThanOneRow
		default:
			return fmt.Errorf("Unhandled Oracle error. Code:[%d] Desc:[%s]", oraError.Code(), oraError.Message())
		}
	} else {
		return err
	}
}
