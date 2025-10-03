package gosqlgen

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

func callingFuncInfo() (int, string) {
	pc, _, n, ok := runtime.Caller(2)
	if !ok {
		return 0, ""
	}
	f := runtime.FuncForPC(pc)
	if f != nil {
		return n, strings.TrimPrefix(f.Name(), "github.com/tsladecek/")
	}
	return n, ""
}

func Errorf(format string, a ...any) error {
	lineNumber, fn := callingFuncInfo()
	format = "\n[%s:%d]:\t" + format
	args := []any{fn, lineNumber}
	args = append(args, a...)
	return fmt.Errorf(format, args...)
}

var (
	ErrFKFieldNumber = errors.New("expected two dot separated fields; expected format: table.column")
	ErrFKTableEmpty  = errors.New("no table specified; expected format: table.column")
	ErrFKColumnEmpty = errors.New("no column specified; expected format: table.column")

	ErrInvalidTagPrefix = errors.New("tag prefix not valid")
	ErrNoClosingQuote   = errors.New("tag not closed with quote")

	ErrEmptyTag          = errors.New("tag empty")
	ErrTagFieldNumber    = errors.New("tag must have at least one field representing column name")
	ErrFKSpecFieldNumber = errors.New("invalid Foreign key spec, must be in format: fk table.column")
	ErrFlagFieldNumber   = errors.New("invalid flag spec")
	ErrFlagFormat        = errors.New("invalid flag format")

	ErrColumnNotFound = errors.New("column not found")

	ErrEmptyTablename = errors.New("tag found in comment group but table name is empty")
	ErrNoTableTag     = errors.New("table tag not found")

	ErrFKTableNotFoundInModel = errors.New("table not found in spec when forming foreign key constraints")

	ErrNoColumnTag = errors.New("no column tag found")

	ErrNoPrimaryKey   = errors.New("no primary key found")
	ErrUnsuportedType = errors.New("unsuported type")

	// valuer
	ErrValuer            = errors.New("failed to infer new value")
	ErrStringKind        = errors.New("unrecognized string kind")
	ErrPrevType          = errors.New("type of previous value does not match valuer")
	ErrValuerConstructor = errors.New("failed to construct valuer")
	ErrValueFormat       = errors.New("failed to format value")
	ErrValueType         = errors.New("invalid type")
)
