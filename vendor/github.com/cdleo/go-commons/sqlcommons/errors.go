package sqlcommons

import (
	"fmt"
)

//Errors
var (
	DBNotSupported   = fmt.Errorf("DataBase not supported")
	DBNotInitialized = fmt.Errorf("DataBase not initialized")
	ConnectionFailed = fmt.Errorf("Connection Failed")
	ConnectionClosed = fmt.Errorf("Connection Closed")
	OpNotSupported   = fmt.Errorf("Operation not supported")
	TxNotFoundInCtx  = fmt.Errorf("SQLTx not found in context")

	UniqueConstraintViolation     = fmt.Errorf("Unique constraint violation")
	IntegrityConstraintViolation  = fmt.Errorf("Integrity constraint violation")
	ValueTooLargeForColumn        = fmt.Errorf("Value too large for column")
	ValueLargerThanPrecision      = fmt.Errorf("Value larger than specified precision")
	CannotSetNullColumn           = fmt.Errorf("Cannot set a NULL value into a NOT NULL column")
	InvalidNumericValue           = fmt.Errorf("Invalid number")
	SubqueryReturnsMoreThanOneRow = fmt.Errorf("Subquery returns more than one row")
)
