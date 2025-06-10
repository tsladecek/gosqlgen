package gosqlgen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
)

type Driver interface {
	Get(w io.Writer, tw io.Writer, table *Table, keys []*Column) error
	Create(w io.Writer, tw io.Writer, table *Table) error
	Update(w io.Writer, tw io.Writer, table *Table, keys []*Column) error
	Delete(w io.Writer, tw io.Writer, table *Table, keys []*Column) error
}

func CreateTemplates(d Driver, model *DBModel) error {
	writer := new(bytes.Buffer)
	testWriter := new(bytes.Buffer)

	writer.Write(fmt.Appendf(nil, `
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
`, model.PackageName))

	testWriter.Write(fmt.Appendf(nil, `
package %s_test
import "testing"
`, model.PackageName))

	for _, table := range model.Tables {
		pk, bk, err := table.PkAndBk()
		if err != nil {
			return fmt.Errorf("Failed to fetch primary and business keys: %w", err)
		}
		err = d.Get(writer, testWriter, table, pk)
		if err != nil {
			return fmt.Errorf("Failed to create GET template by primary keys for table %s: %w", table.Name, err)
		}

		if bk != nil {
			err = d.Get(bufio.NewWriter(writer), testWriter, table, bk)
			if err != nil {
				return fmt.Errorf("Failed to create GET template by business keys for table %s: %w", table.Name, err)
			}
		}
	}

	code, err := format.Source(writer.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to format code: %w", err)
	}
	testCode, err := format.Source(testWriter.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to format test code: %w", err)
	}

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
