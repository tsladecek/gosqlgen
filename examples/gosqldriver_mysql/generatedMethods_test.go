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
		tbl_users_afhjdbje := User{RawId: 1, Id: "aWi3c", Name: []byte(`ax9oly0QA1aLUYnrgpudCjDJAMhprBJd`), payload: []byte(`{"YIXHmd12":"HEI9XuVX", "JN2lXK1O":"F8FZ1boN"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_afhjdbje.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries_ediaegab := Country{RawId: 1, Id: "aVvX1OIKgtQoBUjybDNN68HQ3TVewUnt", Name: "a8VYIilMjgTMbgFix5wcHvmxrIy07STJ", GPS: "azawZmWmh5BZKnZXhpDjaYIzKulGlvts", Continent: "Asia"}
		err = tbl_countries_ediaegab.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_acfiabjk := Address{RawId: 1, Id: "aO6LLoe2pJPGKfbVdK963QHSTyMm6frz", Address: "ae2IeH0xPkFowPZo2x2eu65DQYk9xn9v", UserId: tbl_users_afhjdbje.RawId, CountryId: tbl_countries_ediaegab.RawId}
		err = tbl_addresses_acfiabjk.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Address{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses_acfiabjk.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_acfiabjk, gotByPk)

		// Get By Business Keys
		gotByBk := Address{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses_acfiabjk.Id, tbl_addresses_acfiabjk.Address)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_acfiabjk, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_users_afhjdbje := User{RawId: 1, Id: "aWi3c", Name: []byte(`ax9oly0QA1aLUYnrgpudCjDJAMhprBJd`), payload: []byte(`{"YIXHmd12":"HEI9XuVX", "JN2lXK1O":"F8FZ1boN"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_afhjdbje.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries_ediaegab := Country{RawId: 1, Id: "aVvX1OIKgtQoBUjybDNN68HQ3TVewUnt", Name: "a8VYIilMjgTMbgFix5wcHvmxrIy07STJ", GPS: "azawZmWmh5BZKnZXhpDjaYIzKulGlvts", Continent: "Asia"}
		err = tbl_countries_ediaegab.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_acfiabjk := Address{RawId: 1, Id: "aO6LLoe2pJPGKfbVdK963QHSTyMm6frz", Address: "ae2IeH0xPkFowPZo2x2eu65DQYk9xn9v", UserId: tbl_users_afhjdbje.RawId, CountryId: tbl_countries_ediaegab.RawId}
		err = tbl_addresses_acfiabjk.insert(ctx, testDb)
		require.NoError(t, err)

		err = tbl_addresses_acfiabjk.updateByPrimaryKeys(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Address{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses_acfiabjk.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_acfiabjk, gotByPk)

		// Get By Business Keys

		err = tbl_addresses_acfiabjk.updateByBusinessKeys(ctx, testDb)
		require.NoError(t, err)

		gotByBk := Address{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses_acfiabjk.Id, tbl_addresses_acfiabjk.Address)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_acfiabjk, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_users_afhjdbje := User{RawId: 1, Id: "aWi3c", Name: []byte(`ax9oly0QA1aLUYnrgpudCjDJAMhprBJd`), payload: []byte(`{"YIXHmd12":"HEI9XuVX", "JN2lXK1O":"F8FZ1boN"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_afhjdbje.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries_ediaegab := Country{RawId: 1, Id: "aVvX1OIKgtQoBUjybDNN68HQ3TVewUnt", Name: "a8VYIilMjgTMbgFix5wcHvmxrIy07STJ", GPS: "azawZmWmh5BZKnZXhpDjaYIzKulGlvts", Continent: "Asia"}
		err = tbl_countries_ediaegab.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_acfiabjk := Address{RawId: 1, Id: "aO6LLoe2pJPGKfbVdK963QHSTyMm6frz", Address: "ae2IeH0xPkFowPZo2x2eu65DQYk9xn9v", UserId: tbl_users_afhjdbje.RawId, CountryId: tbl_countries_ediaegab.RawId}
		err = tbl_addresses_acfiabjk.insert(ctx, testDb)
		require.NoError(t, err)

		got := Address{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_addresses_acfiabjk.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_acfiabjk, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := Address{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_addresses_acfiabjk.RawId)
		require.Error(t, err)
	})

}

func TestGoSQLGen_AddressBook(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_users_kbkafhch := User{RawId: 1, Id: "aG0Uq", Name: []byte(`aiuh7Q3Q5yCKhfxgwoEBSG0zy3ztBtbf`), payload: []byte(`{"ZrTrrstG":"9ys8nVup", "FAigTLTi":"VtyV7tRL"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_kbkafhch.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries_kdckklld := Country{RawId: 1, Id: "auLIEDZodKLKQDL3egctszYq10dhIaOw", Name: "a0Rxuy1ORqPSkH7xpISdeQifrDYBS5pD", GPS: "awcRUCm4IKRg6wJ2G3Pecdbq0RfaWswg", Continent: "Asia"}
		err = tbl_countries_kdckklld.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_fdlhacek := Address{RawId: 1, Id: "aJ7G8O5ui46zFdcJZoUgmoV3Md9eOJ6m", Address: "aoc48krMLb9BbiAE9g5F2Ptwx9JJZ9b5", UserId: tbl_users_kbkafhch.RawId, CountryId: tbl_countries_kdckklld.RawId}
		err = tbl_addresses_fdlhacek.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book_hglgacjh := AddressBook{RawId: 1, Id: "aSCfGaCCsZERhzbzc1H99bEeghbpbsMe", AddressId: tbl_addresses_fdlhacek.RawId}
		err = tbl_addresses_book_hglgacjh.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := AddressBook{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses_book_hglgacjh.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book_hglgacjh, gotByPk)

		// Get By Business Keys
		gotByBk := AddressBook{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses_book_hglgacjh.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book_hglgacjh, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_users_kbkafhch := User{RawId: 1, Id: "aG0Uq", Name: []byte(`aiuh7Q3Q5yCKhfxgwoEBSG0zy3ztBtbf`), payload: []byte(`{"ZrTrrstG":"9ys8nVup", "FAigTLTi":"VtyV7tRL"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_kbkafhch.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries_kdckklld := Country{RawId: 1, Id: "auLIEDZodKLKQDL3egctszYq10dhIaOw", Name: "a0Rxuy1ORqPSkH7xpISdeQifrDYBS5pD", GPS: "awcRUCm4IKRg6wJ2G3Pecdbq0RfaWswg", Continent: "Asia"}
		err = tbl_countries_kdckklld.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_fdlhacek := Address{RawId: 1, Id: "aJ7G8O5ui46zFdcJZoUgmoV3Md9eOJ6m", Address: "aoc48krMLb9BbiAE9g5F2Ptwx9JJZ9b5", UserId: tbl_users_kbkafhch.RawId, CountryId: tbl_countries_kdckklld.RawId}
		err = tbl_addresses_fdlhacek.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book_hglgacjh := AddressBook{RawId: 1, Id: "aSCfGaCCsZERhzbzc1H99bEeghbpbsMe", AddressId: tbl_addresses_fdlhacek.RawId}
		err = tbl_addresses_book_hglgacjh.insert(ctx, testDb)
		require.NoError(t, err)

		err = tbl_addresses_book_hglgacjh.updateByPrimaryKeys(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := AddressBook{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_addresses_book_hglgacjh.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book_hglgacjh, gotByPk)

		// Get By Business Keys

		err = tbl_addresses_book_hglgacjh.updateByBusinessKeys(ctx, testDb)
		require.NoError(t, err)

		gotByBk := AddressBook{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_addresses_book_hglgacjh.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book_hglgacjh, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_users_kbkafhch := User{RawId: 1, Id: "aG0Uq", Name: []byte(`aiuh7Q3Q5yCKhfxgwoEBSG0zy3ztBtbf`), payload: []byte(`{"ZrTrrstG":"9ys8nVup", "FAigTLTi":"VtyV7tRL"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_kbkafhch.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_countries_kdckklld := Country{RawId: 1, Id: "auLIEDZodKLKQDL3egctszYq10dhIaOw", Name: "a0Rxuy1ORqPSkH7xpISdeQifrDYBS5pD", GPS: "awcRUCm4IKRg6wJ2G3Pecdbq0RfaWswg", Continent: "Asia"}
		err = tbl_countries_kdckklld.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_fdlhacek := Address{RawId: 1, Id: "aJ7G8O5ui46zFdcJZoUgmoV3Md9eOJ6m", Address: "aoc48krMLb9BbiAE9g5F2Ptwx9JJZ9b5", UserId: tbl_users_kbkafhch.RawId, CountryId: tbl_countries_kdckklld.RawId}
		err = tbl_addresses_fdlhacek.insert(ctx, testDb)
		require.NoError(t, err)
		tbl_addresses_book_hglgacjh := AddressBook{RawId: 1, Id: "aSCfGaCCsZERhzbzc1H99bEeghbpbsMe", AddressId: tbl_addresses_fdlhacek.RawId}
		err = tbl_addresses_book_hglgacjh.insert(ctx, testDb)
		require.NoError(t, err)

		got := AddressBook{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_addresses_book_hglgacjh.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_addresses_book_hglgacjh, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := AddressBook{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_addresses_book_hglgacjh.RawId)
		require.Error(t, err)
	})

}

func TestGoSQLGen_Country(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_countries_djgeigaj := Country{RawId: 1, Id: "auiSyjPLTSkUxnV7QPCVqDOK8ht2XHs1", Name: "aGI5ml8vp9cIUHP4PBd7UodR3tH0wnrQ", GPS: "amyUM4GdcJSJTrkdJvkqV0b9v9RNreyq", Continent: "Asia"}
		err = tbl_countries_djgeigaj.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Country{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_countries_djgeigaj.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries_djgeigaj, gotByPk)

		// Get By Business Keys
		gotByBk := Country{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_countries_djgeigaj.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries_djgeigaj, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_countries_djgeigaj := Country{RawId: 1, Id: "auiSyjPLTSkUxnV7QPCVqDOK8ht2XHs1", Name: "aGI5ml8vp9cIUHP4PBd7UodR3tH0wnrQ", GPS: "amyUM4GdcJSJTrkdJvkqV0b9v9RNreyq", Continent: "Asia"}
		err = tbl_countries_djgeigaj.insert(ctx, testDb)
		require.NoError(t, err)

		tbl_countries_djgeigaj.Name = "bGI5ml8vp9cIUHP4PBd7UodR3tH0wnrQ"
		tbl_countries_djgeigaj.GPS = "bmyUM4GdcJSJTrkdJvkqV0b9v9RNreyq"
		tbl_countries_djgeigaj.Continent = "Europe"
		err = tbl_countries_djgeigaj.updateByPrimaryKeys(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := Country{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_countries_djgeigaj.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries_djgeigaj, gotByPk)

		// Get By Business Keys
		tbl_countries_djgeigaj.Name = "aGI5ml8vp9cIUHP4PBd7UodR3tH0wnrQ"
		tbl_countries_djgeigaj.GPS = "amyUM4GdcJSJTrkdJvkqV0b9v9RNreyq"
		tbl_countries_djgeigaj.Continent = "Asia"
		err = tbl_countries_djgeigaj.updateByBusinessKeys(ctx, testDb)
		require.NoError(t, err)

		gotByBk := Country{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_countries_djgeigaj.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries_djgeigaj, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_countries_djgeigaj := Country{RawId: 1, Id: "auiSyjPLTSkUxnV7QPCVqDOK8ht2XHs1", Name: "aGI5ml8vp9cIUHP4PBd7UodR3tH0wnrQ", GPS: "amyUM4GdcJSJTrkdJvkqV0b9v9RNreyq", Continent: "Asia"}
		err = tbl_countries_djgeigaj.insert(ctx, testDb)
		require.NoError(t, err)

		got := Country{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_countries_djgeigaj.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_countries_djgeigaj, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := Country{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_countries_djgeigaj.RawId)
		require.Error(t, err)
	})

}

func TestGoSQLGen_User(t *testing.T) {
	ctx := t.Context()
	var err error

	t.Run("getInsert", func(t *testing.T) {
		tbl_users_abebbkfj := User{RawId: 1, Id: "a6ntr", Name: []byte(`aYHBwNHwFaUcQogb3y5Uao9EKF9GehpW`), payload: []byte(`{"qOFl1kVz":"LVHXmSje", "tJDK0AMF":"q2Mpxdhn"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_abebbkfj.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := User{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_users_abebbkfj.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_abebbkfj, gotByPk)

		// Get By Business Keys
		gotByBk := User{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_users_abebbkfj.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_abebbkfj, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_users_abebbkfj := User{RawId: 1, Id: "a6ntr", Name: []byte(`aYHBwNHwFaUcQogb3y5Uao9EKF9GehpW`), payload: []byte(`{"qOFl1kVz":"LVHXmSje", "tJDK0AMF":"q2Mpxdhn"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_abebbkfj.insert(ctx, testDb)
		require.NoError(t, err)

		tbl_users_abebbkfj.Name = []byte(`bYHBwNHwFaUcQogb3y5Uao9EKF9GehpW`)
		tbl_users_abebbkfj.payload = []byte(`{"JlPKaKh9":"RuOHeMii", "5AvyjkOI":"kaodvHsf"}`)
		tbl_users_abebbkfj.Age = sql.NullInt32{Valid: true, Int32: 0}
		tbl_users_abebbkfj.DrivesCar = sql.NullBool{Valid: true, Bool: false}
		tbl_users_abebbkfj.Birthday = sql.NullTime{Valid: true, Time: time.Now()}
		tbl_users_abebbkfj.Registered = time.Now()
		err = tbl_users_abebbkfj.updateByPrimaryKeys(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := User{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_users_abebbkfj.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_abebbkfj, gotByPk)

		// Get By Business Keys
		tbl_users_abebbkfj.Name = []byte(`aYHBwNHwFaUcQogb3y5Uao9EKF9GehpW`)
		tbl_users_abebbkfj.payload = []byte(`{"29sTCzpS":"Ll6lY9ij", "PiHwiU9T":"WAsFtnlX"}`)
		tbl_users_abebbkfj.Age = sql.NullInt32{Valid: true, Int32: 1}
		tbl_users_abebbkfj.DrivesCar = sql.NullBool{Valid: true, Bool: true}
		tbl_users_abebbkfj.Birthday = sql.NullTime{Valid: true, Time: time.Now()}
		tbl_users_abebbkfj.Registered = time.Now()
		err = tbl_users_abebbkfj.updateByBusinessKeys(ctx, testDb)
		require.NoError(t, err)

		gotByBk := User{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_users_abebbkfj.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_abebbkfj, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_users_abebbkfj := User{RawId: 1, Id: "a6ntr", Name: []byte(`aYHBwNHwFaUcQogb3y5Uao9EKF9GehpW`), payload: []byte(`{"qOFl1kVz":"LVHXmSje", "tJDK0AMF":"q2Mpxdhn"}`), Age: sql.NullInt32{Valid: true, Int32: 1}, DrivesCar: sql.NullBool{Valid: true, Bool: true}, Birthday: sql.NullTime{Valid: true, Time: time.Now()}, Registered: time.Now()}
		err = tbl_users_abebbkfj.insert(ctx, testDb)
		require.NoError(t, err)

		got := User{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_users_abebbkfj.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_abebbkfj, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := User{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_users_abebbkfj.RawId)
		require.Error(t, err)
	})

}
