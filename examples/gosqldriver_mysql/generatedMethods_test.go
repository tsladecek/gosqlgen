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
