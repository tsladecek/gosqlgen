// This is a generated code by the gosqlgen tool. Do not edit
// see more at: github.com/tsladecek/gosqlgen

package gosqlgen

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
	u.Name = "PFL6K6DQTGRHQRWW7TQX5CVBYD"
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = User{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_users.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// Update By Business Keys
	// Name
	u = gotByBk
	u.Name = "TSXTXNVESHOUUV3F672JJQBTVS"
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
	u.Name = "S6XTLTX6DW5SPT5IDNBG5PTFLT"
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
	u.Name = "FFWGAW6UDHVYVVOXUGIECG7CSM"
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Country{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// GPS
	u = gotByPk
	u.GPS = "EKG6JWX4I7H3ZHXNGVCBT53GDF"
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Country{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.GPS, gotAfterUpdate.GPS)

	// Update By Business Keys
	// Name
	u = gotByBk
	u.Name = "CFSXYPW5R4VPME4VWZPXIVLSZD"
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Country{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_countries.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.Name, gotAfterUpdate.Name)

	// GPS
	u = gotByBk
	u.GPS = "MVC7K7H4S2XC7JJP7EYPPOADMV"
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
	u.UserId = 33
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Address{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.UserId, gotAfterUpdate.UserId)

	// CountryId
	u = gotByPk
	u.CountryId = 71
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Address{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.CountryId, gotAfterUpdate.CountryId)

	// Update By Business Keys
	// UserId
	u = gotByBk
	u.UserId = 91
	err = u.updateByBusinessKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = Address{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses.RawId)
	require.NoError(t, err)
	assert.Equal(t, u.UserId, gotAfterUpdate.UserId)

	// CountryId
	u = gotByBk
	u.CountryId = 123
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
	u.AddressId = 80
	err = u.updateByPrimaryKeys(ctx, testDb)
	require.NoError(t, err)

	gotAfterUpdate = AddressBook{}
	err = gotAfterUpdate.getByPrimaryKeys(ctx, testDb, tbl_addresses_book.RawId)
	require.NoError(t, err)

	assert.Equal(t, u.AddressId, gotAfterUpdate.AddressId)

	// Update By Business Keys
	// AddressId
	u = gotByBk
	u.AddressId = 167
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
