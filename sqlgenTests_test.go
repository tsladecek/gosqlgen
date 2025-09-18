package gosqlgen

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertsAndUpdatedValues(t *testing.T) {
	content, err := format.Source([]byte(strings.Join([]string{
		"package main",
		"// gosqlgen: children",
		"type Child struct {",
		"Id int `gosqlgen:\"id;pk;ai\"`",
		"Name string `gosqlgen:\"name\"`",
		"}",
		"// gosqlgen: parents",
		"type Parent struct {",
		"Id int `gosqlgen:\"id;pk;ai\"`",
		"ChildId int `gosqlgen:\"child_id;fk children.id\"`",
		"}",
	}, "\n")))

	require.NoError(t, err)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	require.NoError(t, err)

	dbModel, err := NewDBModel(fset, f)
	require.NoError(t, err)

	require.NotNil(t, dbModel)
	require.Len(t, dbModel.Tables, 2)

	// inserts
	var childrenTable *Table
	var parentsTable *Table
	for _, t := range dbModel.Tables {
		switch t.Name {
		case "parents":
			parentsTable = t
		case "children":
			childrenTable = t
		}
	}

	require.NotNil(t, childrenTable)
	require.NotNil(t, parentsTable)

	t.Run("children inserts", func(t *testing.T) {
		var w bytes.Buffer
		it, err := childrenTable.testInsert(&w, nil)
		require.NoError(t, err)
		assert.Contains(t, it.varName, "tbl_"+childrenTable.Name)
		assert.Equal(t, childrenTable, it.table)
		assert.Empty(t, it.children)
		assert.Len(t, it.data, 1)

		idCol, err := childrenTable.GetColumn("id")
		require.NoError(t, err)
		require.NotNil(t, idCol)

		nameCol, err := childrenTable.GetColumn("name")
		require.NoError(t, err)
		require.NotNil(t, nameCol)

		assert.Equal(t, nameCol, it.data[0].column)

		valueFormatted, err := it.data[0].value.Format(nameCol.Type)
		require.NoError(t, err)
		expected, err := format.Source(fmt.Appendf(nil, `%s := %s{%s: %s}
err = %s.insert(ctx, testDb)
requireNoError(t, err)
`, it.varName, childrenTable.StructName, nameCol.FieldName, valueFormatted, it.varName))
		require.NoError(t, err)
		actual, err := format.Source(w.Bytes())
		require.NoError(t, err)

		assert.Equal(t, string(expected), string(actual))
	})

	t.Run("parent inserts", func(t *testing.T) {
		var w bytes.Buffer
		it, err := parentsTable.testInsert(&w, nil)
		require.NoError(t, err)
		assert.Contains(t, it.varName, "tbl_"+parentsTable.Name)
		assert.Equal(t, parentsTable, it.table)
		assert.Len(t, it.children, 1)
		require.Len(t, it.data, 0) // FK column values are not saved because their value is injected dynamically

		idCol, err := parentsTable.GetColumn("id")
		require.NoError(t, err)
		require.NotNil(t, idCol)

		childIdCol, err := parentsTable.GetColumn("child_id")
		require.NoError(t, err)
		require.NotNil(t, childIdCol)

		nameCol, err := childrenTable.GetColumn("name")
		require.NoError(t, err)
		require.NotNil(t, nameCol)

		itChild := it.children[0]
		valueFormatted, err := itChild.data[0].value.Format(types.Typ[types.String])
		require.NoError(t, err)

		require.NoError(t, err)
		expected, err := format.Source(fmt.Appendf(nil, `%s := %s{%s: %s}
err = %s.insert(ctx, testDb)
requireNoError(t, err)
%s := %s{%s: %s.%s}
err = %s.insert(ctx, testDb)
requireNoError(t, err)
`, itChild.varName, childrenTable.StructName, nameCol.FieldName, valueFormatted, itChild.varName, it.varName, parentsTable.StructName, childIdCol.FieldName, itChild.varName, childIdCol.ForeignKey.FieldName, it.varName))
		require.NoError(t, err)
		actual, err := format.Source(w.Bytes())
		require.NoError(t, err, w.String())

		assert.Equal(t, string(expected), string(actual))
	})
}
