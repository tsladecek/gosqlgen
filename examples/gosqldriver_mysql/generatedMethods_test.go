// This is a generated code by the gosqlgen tool. Do not edit
// see more at: github.com/tsladecek/gosqlgen

package gosqldrivermysql

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDb *sql.DB

func TestGoSQLGen_Address(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Address{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses, gotByPk)

		// Get By Business Keys
		gotByBk := Address{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses.Id, tbl_addresses.Address)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Address{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses, gotByPk)

		// Get By Business Keys
		gotByBk := Address{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses.Id, tbl_addresses.Address)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)

		got := Address{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := Address{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
		require.Error(t, err)
	})

}

func TestGoSQLGen_AddressBook(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book := AddressBook{AddressId: tbl_addresses.RawId}
		err = tbl_addresses_book.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := AddressBook{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book, gotByPk)

		// Get By Business Keys
		gotByBk := AddressBook{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses_book.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book := AddressBook{AddressId: tbl_addresses.RawId}
		err = tbl_addresses_book.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := AddressBook{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book, gotByPk)

		// Get By Business Keys
		gotByBk := AddressBook{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses_book.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book := AddressBook{AddressId: tbl_addresses.RawId}
		err = tbl_addresses_book.insert(ctx, testDb)
		require.NoError(t, err)

		got := AddressBook{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := AddressBook{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
		require.Error(t, err)
	})

}

func TestGoSQLGen_Country(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Country{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries, gotByPk)

		// Get By Business Keys
		gotByBk := Country{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_countries.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Country{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries, gotByPk)

		// Get By Business Keys
		gotByBk := Country{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_countries.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_countries := Country{}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)

		got := Country{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := Country{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
		require.Error(t, err)
	})

}

func TestGoSQLGen_User(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := User{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users, gotByPk)

		// Get By Business Keys
		gotByBk := User{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_users.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_users, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := User{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users, gotByPk)

		// Get By Business Keys
		gotByBk := User{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_users.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_users, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_users := User{}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)

		got := User{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := User{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
		require.Error(t, err)
	})

}
