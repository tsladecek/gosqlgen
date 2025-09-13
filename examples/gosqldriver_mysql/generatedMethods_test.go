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
		tbl_users_acfiabjk := User{RawId: 1, Id: "aO6LLoe2pJPGKfbVdK963QHSTyMm6frz", Name: []byte(`ae2IeH0xPkFowPZo2x2eu65DQYk9xn9v`)}
		err = tbl_users_acfiabjk.insert(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := User{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_users_afhjdbje.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_afhjdbje, gotByPk)

		// Get By Business Keys
		gotByBk := User{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_users_afhjdbje.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_afhjdbje, gotByBk)
		assert.Equal(t, gotByPk, gotByBk)

	})

	t.Run("update", func(t *testing.T) {
		tbl_users_afhjdbje := User{RawId: 1, Id: "aWi3c0x9oly0QA1aLUYnrgpudCjDJAMh", Name: []byte(`arBJdQglBLIKe1EvHJ1F6QAcW8RxvR2Q`)}
		err = tbl_users_afhjdbje.insert(ctx, testDb)
		require.NoError(t, err)

		tbl_users_afhjdbje.Name = []byte(`brBJdQglBLIKe1EvHJ1F6QAcW8RxvR2Q`)
		err = tbl_users_afhjdbje.updateByPrimaryKeys(ctx, testDb)
		require.NoError(t, err)

		// Get By Primary Keys
		gotByPk := User{}
		err = gotByPk.getByPrimaryKeys(ctx, testDb, tbl_users_afhjdbje.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_afhjdbje, gotByPk)

		// Get By Business Keys
		tbl_users_afhjdbje.Name = []byte(`arBJdQglBLIKe1EvHJ1F6QAcW8RxvR2Q`)
		err = tbl_users_afhjdbje.updateByBusinessKeys(ctx, testDb)
		require.NoError(t, err)

		gotByBk := User{}
		err = gotByBk.getByBusinessKeys(ctx, testDb, tbl_users_afhjdbje.Id)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_afhjdbje, gotByBk)

	})

	t.Run("delete", func(t *testing.T) {
		tbl_users_debcdkil := User{RawId: 1, Id: "amd12HEI9XuVXJN2lXK1OF8FZ1boNKVg", Name: []byte(`akY8RZVvX1OIKgtQoBUjybDNN68HQ3TV`)}
		err = tbl_users_debcdkil.insert(ctx, testDb)
		require.NoError(t, err)

		got := User{}
		err = got.getByPrimaryKeys(ctx, testDb, tbl_users_afhjdbje.RawId)
		require.NoError(t, err)
		assert.Equal(t, tbl_users_afhjdbje, got)

		err = got.delete(ctx, testDb)
		require.NoError(t, err)
		gotAfterDelete := User{}
		err = gotAfterDelete.getByPrimaryKeys(ctx, testDb, tbl_users_afhjdbje.RawId)
		require.Error(t, err)
	})

}
