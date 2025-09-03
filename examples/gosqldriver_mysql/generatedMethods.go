// This is a generated code by the gosqlgen tool. Do not edit
// see more at: github.com/tsladecek/gosqlgen

package gosqldrivermysql

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

func (t *User) getByPrimaryKeys(ctx context.Context, db dbExecutor, _id int) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, name FROM users WHERE _id = ?", _id).Scan(&t.RawId, &t.Id, &t.Name)

	if err != nil {
		return err
	}

	return nil
}

func (t *User) getByBusinessKeys(ctx context.Context, db dbExecutor, id string) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, name FROM users WHERE id = ?", id).Scan(&t.RawId, &t.Id, &t.Name)

	if err != nil {
		return err
	}

	return nil
}

func (t *User) insert(ctx context.Context, db dbExecutor) error {
	res, err := db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?)", t.Id, t.Name)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.RawId = int(id)

	return nil
}

func (t *User) updateByPrimaryKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE users SET name = ? WHERE _id=?", t.Name, t.RawId)
	return err
}

func (t *User) updateByBusinessKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE users SET name = ? WHERE id=?", t.Name, t.Id)
	return err
}

func (t *User) delete(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "DELETE FROM users WHERE _id=?", t.RawId)
	return err
}
