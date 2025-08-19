package gosqlgen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
)

type Driver interface {
	Get(w io.Writer, table *Table, keys []*Column, methodName string) error
	Create(w io.Writer, table *Table, methodName string) error
	Update(w io.Writer, table *Table, keys []*Column, methodName string) error
	Delete(w io.Writer, table *Table, keys []*Column, methodName string) error
}

type MethodName string

const (
	MethodGetByPrimaryKeys     MethodName = "getByPrimaryKeys"
	MethodGetByBusinessKeys    MethodName = "getByBusinessKeys"
	MethodInsert               MethodName = "insert"
	MethodUpdateByPrimaryKeys  MethodName = "updateByPrimaryKeys"
	MethodUpdateByBusinessKeys MethodName = "updateByBusinessKeys"
	MethodDelete               MethodName = "delete"
)

const DBExecutorVarName = "testSqlDb"

func CreateTemplates(d Driver, model *DBModel) error {
	writer := new(bytes.Buffer)
	testWriter := new(bytes.Buffer)

	header := `// This is a generated code by the gosqlgen tool. Do not edit
// see more at: github.com/tsladecek/gosqlgen
`

	writer.Write(fmt.Appendf(nil, `
%s

package %s
import (
	"context"
	"database/sql"
)
type dbExecutor interface {
	// ExecContext executes a query without returning any rows. The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// PrepareContext creates a prepared statement for later queries or executions. Multiple queries or executions may be run concurrently from the returned statement. The caller must call the statement's *Stmt.Close method when the statement is no longer needed.
	// The provided context is used for the preparation of the statement, not for the execution of the statement.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// QueryContext executes a query that returns rows, typically a SELECT. The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// QueryRowContext executes a query that is expected to return at most one row. QueryRowContext always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row.Scan will return ErrNoRows. Otherwise, *Row.Scan scans the first selected row and discards the rest.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
`, header, model.PackageName))

	testWriter.Write(fmt.Appendf(nil, `
%s

package %s
import (
	"testing"
	"database/sql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDb *sql.DB
`, header, model.PackageName))

	ts, err := NewTestSuite()
	if err != nil {
		return err
	}

	for _, table := range model.Tables {
		// GET
		pk, bk, err := table.PkAndBk()
		if err != nil {
			return fmt.Errorf("Failed to fetch primary and business keys: %w", err)
		}
		err = d.Get(writer, table, pk, string(MethodGetByPrimaryKeys))
		if err != nil {
			return fmt.Errorf("Failed to create GET template by primary keys for table %s: %w", table.Name, err)
		}

		if len(bk) > 0 {
			err = d.Get(writer, table, bk, string(MethodGetByBusinessKeys))
			if err != nil {
				return fmt.Errorf("Failed to create GET template by business keys for table %s: %w", table.Name, err)
			}
		}

		// CREATE
		err = d.Create(writer, table, string(MethodInsert))
		if err != nil {
			return fmt.Errorf("Failed to create insert template for table %s: %w", table.Name, err)
		}

		// UPDATE
		err = d.Update(writer, table, pk, string(MethodUpdateByPrimaryKeys))
		if err != nil {
			return fmt.Errorf("Failed to create update template for table %s by primary keys: %w", table.Name, err)
		}

		if len(bk) > 0 {
			err = d.Update(writer, table, bk, string(MethodUpdateByBusinessKeys))
			if err != nil {
				return fmt.Errorf("Failed to create update template for table %s by business keys: %w", table.Name, err)
			}
		}

		// DELETE
		err = d.Delete(writer, table, pk, string(MethodDelete))
		if err != nil {
			return fmt.Errorf("Failed to create delete template for table %s by primary keys: %w", table.Name, err)
		}

		err = ts.Generate(testWriter, table)
		if err != nil {
			return fmt.Errorf("Failed to create TestGET template by primary keys for table %s: %w", table.Name, err)
		}
	}

	code, err := format.Source(writer.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to format code: %w %v", err, writer.String())
	}
	// code := writer.Bytes()

	testCode, err := format.Source(testWriter.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to format test code: %w", err)
	}

	// testCode := testWriter.Bytes()

	err = os.WriteFile("generatedMethods.go", code, 0666)
	if err != nil {
		return fmt.Errorf("Failed writing code to a file: %w", err)
	}

	err = os.WriteFile("generatedMethods_test.go", testCode, 0666)
	if err != nil {
		return fmt.Errorf("Failed writing test code to a file: %w", err)
	}

	return nil
}
