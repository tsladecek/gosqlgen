// This is a generated code by the gosqlgen tool. Do not edit
// see more at: github.com/tsladecek/gosqlgen

package gosqlgen

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

func (t *Admin) getByPrimaryKeys(ctx context.Context, db dbExecutor, _id int) error {
	err := db.QueryRowContext(ctx, "SELECT _id, name FROM admins WHERE _id = ?", _id).Scan(&t.RawId, &t.Name)

	if err != nil {
		return err
	}

	return nil
}

func (t *Admin) insert(ctx context.Context, db dbExecutor) error {
	res, err := db.ExecContext(ctx, "INSERT INTO admins (name) VALUES (?)", t.Name)
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

func (t *Admin) updateByPrimaryKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE admins SET name = ? WHERE _id=?", t.Name, t.RawId)
	return err
}

func (t *Admin) delete(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "DELETE FROM admins WHERE _id=?", t.RawId)
	return err
}

func (t *Country) getByPrimaryKeys(ctx context.Context, db dbExecutor, _id int) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, name, gps FROM countries WHERE _id = ?", _id).Scan(&t.RawId, &t.Id, &t.Name, &t.GPS)

	if err != nil {
		return err
	}

	return nil
}

func (t *Country) getByBusinessKeys(ctx context.Context, db dbExecutor, id string) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, name, gps FROM countries WHERE id = ?", id).Scan(&t.RawId, &t.Id, &t.Name, &t.GPS)

	if err != nil {
		return err
	}

	return nil
}

func (t *Country) insert(ctx context.Context, db dbExecutor) error {
	res, err := db.ExecContext(ctx, "INSERT INTO countries (id, name, gps) VALUES (?, ?, ?)", t.Id, t.Name, t.GPS)
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

func (t *Country) updateByPrimaryKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE countries SET name = ?, gps = ? WHERE _id=?", t.Name, t.GPS, t.RawId)
	return err
}

func (t *Country) updateByBusinessKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE countries SET name = ?, gps = ? WHERE id=?", t.Name, t.GPS, t.Id)
	return err
}

func (t *Country) delete(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "DELETE FROM countries WHERE _id=?", t.RawId)
	return err
}

func (t *Address) getByPrimaryKeys(ctx context.Context, db dbExecutor, _id int32) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, address, user_id, country_id, deleted_at FROM addresses WHERE _id = ? AND deleted_at IS NOT NULL", _id).Scan(&t.RawId, &t.Id, &t.Address, &t.UserId, &t.CountryId, &t.DeletedAt)

	if err != nil {
		return err
	}

	return nil
}

func (t *Address) getByBusinessKeys(ctx context.Context, db dbExecutor, id string, address string) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, address, user_id, country_id, deleted_at FROM addresses WHERE id = ? AND address = ? AND deleted_at IS NOT NULL", id, address).Scan(&t.RawId, &t.Id, &t.Address, &t.UserId, &t.CountryId, &t.DeletedAt)

	if err != nil {
		return err
	}

	return nil
}

func (t *Address) insert(ctx context.Context, db dbExecutor) error {
	res, err := db.ExecContext(ctx, "INSERT INTO addresses (id, address, user_id, country_id) VALUES (?, ?, ?, ?)", t.Id, t.Address, t.UserId, t.CountryId)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.RawId = int32(id)

	return nil
}

func (t *Address) updateByPrimaryKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE addresses SET user_id = ?, country_id = ? WHERE _id=?", t.UserId, t.CountryId, t.RawId)
	return err
}

func (t *Address) updateByBusinessKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE addresses SET user_id = ?, country_id = ? WHERE id=? AND address=?", t.UserId, t.CountryId, t.Id, t.Address)
	return err
}

func (t *Address) delete(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE addresses SET deleted_at = CURRENT_TIMESTAMP WHERE _id=?", t.RawId)
	return err
}

func (t *AddressBook) getByPrimaryKeys(ctx context.Context, db dbExecutor, _id int) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, address_id FROM addresses_book WHERE _id = ?", _id).Scan(&t.RawId, &t.Id, &t.AddressId)

	if err != nil {
		return err
	}

	return nil
}

func (t *AddressBook) getByBusinessKeys(ctx context.Context, db dbExecutor, id string) error {
	err := db.QueryRowContext(ctx, "SELECT _id, id, address_id FROM addresses_book WHERE id = ?", id).Scan(&t.RawId, &t.Id, &t.AddressId)

	if err != nil {
		return err
	}

	return nil
}

func (t *AddressBook) insert(ctx context.Context, db dbExecutor) error {
	res, err := db.ExecContext(ctx, "INSERT INTO addresses_book (id, address_id) VALUES (?, ?)", t.Id, t.AddressId)
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

func (t *AddressBook) updateByPrimaryKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE addresses_book SET address_id = ? WHERE _id=?", t.AddressId, t.RawId)
	return err
}

func (t *AddressBook) updateByBusinessKeys(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE addresses_book SET address_id = ? WHERE id=?", t.AddressId, t.Id)
	return err
}

func (t *AddressBook) delete(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "DELETE FROM addresses_book WHERE _id=?", t.RawId)
	return err
}
