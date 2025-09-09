// This is a generated code by the gosqlgen tool. Do not edit
// see more at: github.com/tsladecek/gosqlgen

package gosqldrivermysql

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDb *sql.DB

func TestGoSQLGen_Address(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_users := User{RawId: 1, Id: "ahYwY", Name: []byte(`aNBnGtXR6V1ZYDT34K4SL1hwzkwH82MP`), payload: []byte(`{"NPHDjfJt":"wNB9xfNw", "C1c4af6l":"tRgNrcq3"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{RawId: 1, Id: "aI3RYdgHL6YPvdFMVgao4naDbqJims9r", Name: "ah5jUmM9PRK2aaqdgA9zQtjA644r10yf", GPS: "ahAOFgs87I9krI1mWDilh7wUdGe16ncw", Continent: "Asia"}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{RawId: 1, Id: "aPUPeTmgkdtzh7qOnPunvMWzyfU3FTYJ", Address: "a6mQbrVeElBGKc8a4j3qPQ0Smc06J7fp", UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
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
		tbl_users := User{RawId: 1, Id: "ahYwY", Name: []byte(`aNBnGtXR6V1ZYDT34K4SL1hwzkwH82MP`), payload: []byte(`{"NPHDjfJt":"wNB9xfNw", "C1c4af6l":"tRgNrcq3"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{RawId: 1, Id: "aI3RYdgHL6YPvdFMVgao4naDbqJims9r", Name: "ah5jUmM9PRK2aaqdgA9zQtjA644r10yf", GPS: "ahAOFgs87I9krI1mWDilh7wUdGe16ncw", Continent: "Asia"}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{RawId: 1, Id: "aPUPeTmgkdtzh7qOnPunvMWzyfU3FTYJ", Address: "a6mQbrVeElBGKc8a4j3qPQ0Smc06J7fp", UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
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
		tbl_users := User{RawId: 1, Id: "ahYwY", Name: []byte(`aNBnGtXR6V1ZYDT34K4SL1hwzkwH82MP`), payload: []byte(`{"NPHDjfJt":"wNB9xfNw", "C1c4af6l":"tRgNrcq3"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{RawId: 1, Id: "aI3RYdgHL6YPvdFMVgao4naDbqJims9r", Name: "ah5jUmM9PRK2aaqdgA9zQtjA644r10yf", GPS: "ahAOFgs87I9krI1mWDilh7wUdGe16ncw", Continent: "Asia"}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{RawId: 1, Id: "aPUPeTmgkdtzh7qOnPunvMWzyfU3FTYJ", Address: "a6mQbrVeElBGKc8a4j3qPQ0Smc06J7fp", UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
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
		tbl_users := User{RawId: 1, Id: "arMNA", Name: []byte(`ajRH9C50Os6KilxPpASu0LElmOKF9ML6`), payload: []byte(`{"HB39I1Um":"aJIHqoJI", "X3814g5u":"DEzf9h9R"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{RawId: 1, Id: "afUXF68PkeoRCfGSAz5E8EHOS84uJFCY", Name: "a57ljC1L1uMjV5lzZBhIiGltMrqqVrVw", GPS: "a8zDCfT1dGtCR8bD8t5UTcEgkCQjaAFQ", Continent: "Asia"}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{RawId: 1, Id: "afjcKazm8bt4d2pEhgyDvSv3LYa5hEZ2", Address: "a6ZXPByWkwsSYfXJPowmO8pV7xAXgNLC", UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book := AddressBook{RawId: 1, Id: "abpbZ7QIguAu4C1YuIRGRKKE6FWv9JEY", AddressId: tbl_addresses.RawId}
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
		tbl_users := User{RawId: 1, Id: "arMNA", Name: []byte(`ajRH9C50Os6KilxPpASu0LElmOKF9ML6`), payload: []byte(`{"HB39I1Um":"aJIHqoJI", "X3814g5u":"DEzf9h9R"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{RawId: 1, Id: "afUXF68PkeoRCfGSAz5E8EHOS84uJFCY", Name: "a57ljC1L1uMjV5lzZBhIiGltMrqqVrVw", GPS: "a8zDCfT1dGtCR8bD8t5UTcEgkCQjaAFQ", Continent: "Asia"}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{RawId: 1, Id: "afjcKazm8bt4d2pEhgyDvSv3LYa5hEZ2", Address: "a6ZXPByWkwsSYfXJPowmO8pV7xAXgNLC", UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book := AddressBook{RawId: 1, Id: "abpbZ7QIguAu4C1YuIRGRKKE6FWv9JEY", AddressId: tbl_addresses.RawId}
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
		tbl_users := User{RawId: 1, Id: "arMNA", Name: []byte(`ajRH9C50Os6KilxPpASu0LElmOKF9ML6`), payload: []byte(`{"HB39I1Um":"aJIHqoJI", "X3814g5u":"DEzf9h9R"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
		err = tbl_users.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries := Country{RawId: 1, Id: "afUXF68PkeoRCfGSAz5E8EHOS84uJFCY", Name: "a57ljC1L1uMjV5lzZBhIiGltMrqqVrVw", GPS: "a8zDCfT1dGtCR8bD8t5UTcEgkCQjaAFQ", Continent: "Asia"}
		err = tbl_countries.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses := Address{RawId: 1, Id: "afjcKazm8bt4d2pEhgyDvSv3LYa5hEZ2", Address: "a6ZXPByWkwsSYfXJPowmO8pV7xAXgNLC", UserId: tbl_users.RawId, CountryId: tbl_countries.RawId}
		err = tbl_addresses.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book := AddressBook{RawId: 1, Id: "abpbZ7QIguAu4C1YuIRGRKKE6FWv9JEY", AddressId: tbl_addresses.RawId}
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
		tbl_countries := Country{RawId: 1, Id: "aoC9xBCendi79bdsYS4YkYf7pVjKvEmh", Name: "a6EY2UPfMYHCL570ZXGdJg1jaRCX1ujF", GPS: "agGcjxIQlA2LKfMLklWJRczY0zWIlTXy", Continent: "Asia"}
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
		tbl_countries := Country{RawId: 1, Id: "aoC9xBCendi79bdsYS4YkYf7pVjKvEmh", Name: "a6EY2UPfMYHCL570ZXGdJg1jaRCX1ujF", GPS: "agGcjxIQlA2LKfMLklWJRczY0zWIlTXy", Continent: "Asia"}
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
		tbl_countries := Country{RawId: 1, Id: "aoC9xBCendi79bdsYS4YkYf7pVjKvEmh", Name: "a6EY2UPfMYHCL570ZXGdJg1jaRCX1ujF", GPS: "agGcjxIQlA2LKfMLklWJRczY0zWIlTXy", Continent: "Asia"}
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
		tbl_users := User{RawId: 1, Id: "a9fQs", Name: []byte(`auuZt5Awh5I7ZkT1IO0p2NCYxh4lbsJQ`), payload: []byte(`{"sPrJ477G":"AQzqrPgI", "j5snJ9Av":"P2rWKsD9"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
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
		tbl_users := User{RawId: 1, Id: "a9fQs", Name: []byte(`auuZt5Awh5I7ZkT1IO0p2NCYxh4lbsJQ`), payload: []byte(`{"sPrJ477G":"AQzqrPgI", "j5snJ9Av":"P2rWKsD9"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
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
		tbl_users := User{RawId: 1, Id: "a9fQs", Name: []byte(`auuZt5Awh5I7ZkT1IO0p2NCYxh4lbsJQ`), payload: []byte(`{"sPrJ477G":"AQzqrPgI", "j5snJ9Av":"P2rWKsD9"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}}
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
