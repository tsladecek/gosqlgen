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

func TestGoSQLGen_User(t *testing.T) {
	ctx := t.Context()
	var err error

	// Inserts
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

	var gotAfterUpdate User
	var u User

	// Update By Primary Keys
	// Name
	u = gotByPk
	u.Name = "OOC23MA73GNGP63TGFID5GBXLE"
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = User{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// Update By Business Keys
	// Name
	u = gotByBk
	u.Name = "LMXXHSC5LPLFWUN4KEMJINPWLJ"
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = User{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// Delete
	err = gotByPk.delete(ctx, testDb)
	require.NoError(t, err)
	gotAfterDelete := User{}
	err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
	require.Error(t, err)
}

func TestGoSQLGen_Admin(t *testing.T) {
	ctx := t.Context()
	var err error

	// Inserts
	tbl_users := User{}
	err = tbl_users.insert(ctx, testDb)
	require.NoError(t, err)

	tbl_admins := Admin{RawId: tbl_users.RawId}
	err = tbl_admins.insert(ctx, testDb)
	require.NoError(t, err)

	// Get By Primary Keys
	gotByPk := Admin{}
	err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_admins.RawId)
	require.NoError(t, err)
	assert.Equal(t, tbl_admins, gotByPk)

	var gotAfterUpdate Admin
	var u Admin

	// Update By Primary Keys
	// Name
	u = gotByPk
	u.Name = "7KY2P3NL7LT2PEXTPYHZJZ2ZQI"
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Admin{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_admins.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// Delete
	err = gotByPk.delete(ctx, testDb)
	require.NoError(t, err)
	gotAfterDelete := Admin{}
	err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_admins.RawId)
	require.Error(t, err)
}

func TestGoSQLGen_Country(t *testing.T) {
	ctx := t.Context()
	var err error

	// Inserts
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

	var gotAfterUpdate Country
	var u Country

	// Update By Primary Keys
	// Name
	u = gotByPk
	u.Name = "AJFILI4JKZQRXVDJMGRZUG5LIJ"
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Country{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// GPS
	u = gotByPk
	u.GPS = "6D6PR5VK7OCZK3ANBVQGQBY2ZH"
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Country{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.GPS, gotAfterUpdate.GPS)

	// Update By Business Keys
	// Name
	u = gotByBk
	u.Name = "IBLHDWG55HI4O6R5HRBQR6NPUJ"
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Country{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// GPS
	u = gotByBk
	u.GPS = "LTTS4TJNXMVYBOTV25GV4SM7HB"
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Country{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.GPS, gotAfterUpdate.GPS)

	// Delete
	err = gotByPk.delete(ctx, testDb)
	require.NoError(t, err)
	gotAfterDelete := Country{}
	err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.Error(t, err)
}

func TestGoSQLGen_Address(t *testing.T) {
	ctx := t.Context()
	var err error

	// Inserts
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

	var gotAfterUpdate Address
	var u Address

	// Update By Primary Keys
	// UserId
	u = gotByPk
	u.UserId = 63
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Address{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.UserId, gotAfterUpdate.UserId)

	// CountryId
	u = gotByPk
	u.CountryId = 189
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Address{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.CountryId, gotAfterUpdate.CountryId)

	// Update By Business Keys
	// UserId
	u = gotByBk
	u.UserId = 244
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Address{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.UserId, gotAfterUpdate.UserId)

	// CountryId
	u = gotByBk
	u.CountryId = 116
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Address{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.CountryId, gotAfterUpdate.CountryId)

	// Delete
	err = gotByPk.delete(ctx, testDb)
	require.NoError(t, err)
	gotAfterDelete := Address{}
	err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.Error(t, err)
}

func TestGoSQLGen_AddressBook(t *testing.T) {
	ctx := t.Context()
	var err error

	// Inserts
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

	var gotAfterUpdate AddressBook
	var u AddressBook

	// Update By Primary Keys
	// AddressId
	u = gotByPk
	u.AddressId = 41
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = AddressBook{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.AddressId, gotAfterUpdate.AddressId)

	// Update By Business Keys
	// AddressId
	u = gotByBk
	u.AddressId = 170
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = AddressBook{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.AddressId, gotAfterUpdate.AddressId)

	// Delete
	err = gotByPk.delete(ctx, testDb)
	require.NoError(t, err)
	gotAfterDelete := AddressBook{}
	err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
	require.Error(t, err)
}
