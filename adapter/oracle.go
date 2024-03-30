package adapter

import (
	"fmt"

	"github.com/cdleo/go-commons/sqlcommons"
	"github.com/godror/godror"
)

type oracleAdapter struct{}

func NewOracleAdapter() sqlcommons.SQLAdapter {
	return &oracleAdapter{}
}

func (t *oracleAdapter) Translate(query string) string {
	return query
}

func (s *oracleAdapter) ErrorHandler(err error) error {
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
