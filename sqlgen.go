package gosqlgen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
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

func additionalImports(model *DBModel) ([]string, []string, error) {
	codeImports := []string{}
	testCodeImports := []string{}
	testCodeImportsMap := make(map[string]bool)
	for _, table := range model.Tables {
		if table.SkipTests {
			continue
		}

		for _, column := range table.Columns {
			if IsOneOfTypes(column.Type, []string{"time.Time"}) {
				testCodeImportsMap["time"] = true
			}
		}
	}

	if len(testCodeImportsMap) > 0 {
		testCodeImports = slices.Collect(maps.Keys(testCodeImportsMap))
	}

	return codeImports, testCodeImports, nil
}

func formatImports(imports []string) string {
	formatted := []string{}
	for _, imp := range imports {
		formatted = append(formatted, fmt.Sprintf(`"%s"`, imp))
	}
	return strings.Join(formatted, "\n")
}

func CreateTemplates(d Driver, model *DBModel, outputPath, outputTestPath string) error {
	writer := new(bytes.Buffer)
	writerContent := new(bytes.Buffer)
	testWriter := new(bytes.Buffer)
	testWriterContent := new(bytes.Buffer)

	header := `// This is a generated code by the gosqlgen tool. Do not edit
// see more at: github.com/tsladecek/gosqlgen
`

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
		err = d.Get(writerContent, table, pk, string(MethodGetByPrimaryKeys))
		if err != nil {
			return fmt.Errorf("Failed to create GET template by primary keys for table %s: %w", table.Name, err)
		}

		if len(bk) > 0 {
			err = d.Get(writerContent, table, bk, string(MethodGetByBusinessKeys))
			if err != nil {
				return fmt.Errorf("Failed to create GET template by business keys for table %s: %w", table.Name, err)
			}
		}

		// CREATE
		err = d.Create(writerContent, table, string(MethodInsert))
		if err != nil {
			return fmt.Errorf("Failed to create insert template for table %s: %w", table.Name, err)
		}

		// UPDATE
		err = d.Update(writerContent, table, pk, string(MethodUpdateByPrimaryKeys))
		if err != nil {
			return fmt.Errorf("Failed to create update template for table %s by primary keys: %w", table.Name, err)
		}

		if len(bk) > 0 {
			err = d.Update(writerContent, table, bk, string(MethodUpdateByBusinessKeys))
			if err != nil {
				return fmt.Errorf("Failed to create update template for table %s by business keys: %w", table.Name, err)
			}
		}

		// DELETE
		err = d.Delete(writerContent, table, pk, string(MethodDelete))
		if err != nil {
			return fmt.Errorf("Failed to create delete template for table %s by primary keys: %w", table.Name, err)
		}

		if !table.SkipTests {
			err = ts.Generate(testWriterContent, table)
			if err != nil {
				return fmt.Errorf("Failed to create test template for table %s: %w", table.Name, err)
			}
		}
	}

	codeImportsRaw, testCodeImportsRaw, err := additionalImports(model)
	if err != nil {
		return fmt.Errorf("%w: when inferring additional imports", err)
	}
	codeImports := formatImports(codeImportsRaw)
	testCodeImports := formatImports(testCodeImportsRaw)

	writer.Write(fmt.Appendf(nil, `
%s

package %s
import (
	"context"
	"database/sql"
	%s
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
`, header, model.PackageName, codeImports))

	testWriter.Write(fmt.Appendf(nil, `
%s

package %s
import (
	"testing"
	"database/sql"
	%s

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDb *sql.DB
`, header, model.PackageName, testCodeImports))

	_, err = io.Copy(writer, writerContent)
	if err != nil {
		return fmt.Errorf("%w: when writing content", err)
	}

	_, err = io.Copy(testWriter, testWriterContent)
	if err != nil {
		return fmt.Errorf("%w: when writing test content", err)
	}

	code, err := format.Source(writer.Bytes())
	if err != nil {
		return fmt.Errorf("%w: when formating code", err)
	}

	testCode, err := format.Source(testWriter.Bytes())
	if err != nil {
		return fmt.Errorf("%w: when formating test code", err)
	}

	err = os.WriteFile(outputPath, code, 0666)
	if err != nil {
		return fmt.Errorf("%w: when writing code to a file", err)
	}

	err = os.WriteFile(outputTestPath, testCode, 0666)
	if err != nil {
		return fmt.Errorf("%w: when writing test code to a file", err)
	}

	return nil
}
