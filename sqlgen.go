package gosqlgen

import (
	"fmt"
	"io"
	"os"
)

type Driver interface {
	Get(w io.Writer, tw io.Writer, table *Table, keys []*Column) error
	Create(w io.Writer, tw io.Writer, table *Table) error
	Update(w io.Writer, tw io.Writer, table *Table, keys []*Column) error
	Delete(w io.Writer, tw io.Writer, table *Table, keys []*Column) error
}

func SaveTemplate(t string, path string) error {
	fmt.Println(t)
	return nil
}

func CreateTemplates(d Driver, model *DBModel) error {
	writer := os.Stdout
	testWriter := os.Stdout

	writer.Write([]byte(fmt.Sprintf(`package %s

type DbExecutor interface {
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
`, model.PackageName)))

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
			err = d.Get(writer, testWriter, table, bk)
			if err != nil {
				return fmt.Errorf("Failed to create GET template by business keys for table %s: %w", table.Name, err)
			}
		}
	}

	return nil
}
