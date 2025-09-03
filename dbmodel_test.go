package gosqlgen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFKTableAndColumn(t *testing.T) {
	cases := []struct {
		name           string
		fk             string
		expectedTable  string
		expectedColumn string
		expectedErr    error
	}{
		{name: "valid", fk: "table.column", expectedErr: nil, expectedTable: "table", expectedColumn: "column"},
		{name: "invalid - more than one field separated by dot", fk: "table.column.column", expectedErr: ErrFKFieldNumber},
		{name: "invalid - less than one field separated by dot", fk: "table", expectedErr: ErrFKFieldNumber},
		{name: "invalid - empty string", fk: "", expectedErr: ErrFKFieldNumber},
		{name: "invalid - table empty", fk: ".column", expectedErr: ErrFKTableEmpty},
		{name: "invalid - column empty", fk: "table.", expectedErr: ErrFKColumnEmpty},
		{name: "invalid - table space", fk: "  .column", expectedErr: ErrFKTableEmpty},
		{name: "invalid - column space", fk: "table.  ", expectedErr: ErrFKColumnEmpty},
		{name: "valid - trim spaces", fk: "  table  .  column  ", expectedErr: nil, expectedTable: "table", expectedColumn: "column"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := Column{fk: tt.fk}
			ta, co, err := c.FKTableAndColumn()
			require.Equal(t, tt.expectedErr == nil, err == nil)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedTable, ta)
				assert.Equal(t, tt.expectedColumn, co)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestExtractTagContent(t *testing.T) {
	cases := []struct {
		name        string
		tagName     string
		input       string
		output      string
		expectedErr error
	}{
		{name: "valid", tagName: "tag", input: `tag:"input"`, output: "input", expectedErr: nil},
		{name: "invalid - missing tag prefix in input", tagName: "tag", input: `:"input"`, expectedErr: ErrInvalidTagPrefix},
		{name: "invalid - missing colon", tagName: "tag", input: `tag"input"`, expectedErr: ErrInvalidTagPrefix},
		{name: "invalid - double quote not exactly after colon", tagName: "tag", input: `tag:tag"input"`, expectedErr: ErrInvalidTagPrefix},
		{name: "invalid - missing start quote", tagName: "tag", input: `tag:input"`, expectedErr: ErrInvalidTagPrefix},
		{name: "invalid - missing end quote", tagName: "tag", input: `tag:"input`, expectedErr: ErrNoClosingQuote},
		{name: "valid - output should be space trimmed", tagName: "tag", input: `tag:"  input "`, output: "input", expectedErr: nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ExtractTagContent(tt.tagName, tt.input)
			require.Equal(t, tt.expectedErr == nil, err == nil)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.output, c)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestNewColumn(t *testing.T) {
	cases := []struct {
		name           string
		tag            string
		expectedErr    error
		expectedColumn Column
	}{
		{name: "invalid - empty tag", tag: fmt.Sprintf(`%s:""`, TagPrefix), expectedErr: ErrEmptyTag},
		{name: "invalid - tag parsing", tag: fmt.Sprintf(`%s:col`, TagPrefix), expectedErr: ErrInvalidTagPrefix},
		{name: "invalid - fk spec contains more than two space separated fields", tag: fmt.Sprintf(`%s:"column;int;pk ai;fk table col;bk;sd"`, TagPrefix), expectedErr: ErrFKSpecFieldNumber},
		{name: "invalid - fk spec contains less than two space separated fields", tag: fmt.Sprintf(`%s:"column;int;pk ai;fk;bk;sd"`, TagPrefix), expectedErr: ErrFKSpecFieldNumber},
		{name: "valid - column with everything", tag: fmt.Sprintf(`%s:"column;int;pk;ai;fk table.col;bk;sd"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true, SoftDelete: true, BusinessKey: true, AutoIncrement: true, fk: "table.col"}},
		{name: "valid - column with everything with spaces that should be trimmed", tag: fmt.Sprintf(`%s:"   column  ;  int   ;    pk;ai   ;     fk table.col   ;  bk  ;  sd  "`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true, SoftDelete: true, BusinessKey: true, AutoIncrement: true, fk: "table.col"}},
		{name: "valid - just pk", tag: fmt.Sprintf(`%s:"column;int;pk"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true}},
		{name: "valid - unrecognized tag is skipped", tag: fmt.Sprintf(`%s:"column;int;bad"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column"}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewColumn(tt.tag)
			require.Equal(t, tt.expectedErr == nil, err == nil)
			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedColumn, *c)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestGetColumn(t *testing.T) {
	col1 := &Column{Name: "col1"}
	col2 := &Column{Name: "col2"}
	cases := []struct {
		name        string
		expectedErr error
		table       *Table
		columnName  string
		expectedCol *Column
	}{
		{name: "valid - col1 found", table: &Table{Columns: []*Column{col1, col2}}, columnName: "col1", expectedCol: col1},
		{name: "valid - col2 found", table: &Table{Columns: []*Column{col1, col2}}, columnName: "col2", expectedCol: col2},
		{name: "invalid - col3 not found", table: &Table{Columns: []*Column{col1, col2}}, columnName: "col3", expectedErr: ErrColumnNotFound},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := tt.table.GetColumn(tt.columnName)
			require.Equal(t, tt.expectedErr == nil, err == nil)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedCol, c)
			} else {
				assert.ErrorIs(t, err, ErrColumnNotFound)
			}
		})
	}
}

func TestParseTableName(t *testing.T) {
	cases := []struct {
		name          string
		expectedErr   error
		expectedTable Table
		comments      []string
	}{
		{name: "invalid no annotation found", comments: []string{"// comment line 1", "// comment line 2"}, expectedErr: ErrNoTableTag},
		{name: "invalid empty table name", comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s:", TagPrefix)}, expectedErr: ErrEmptyTablename},
		{name: "invalid empty table name of spaces", comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s:   ", TagPrefix)}, expectedErr: ErrEmptyTablename},
		{name: "valid", expectedTable: Table{Name: "table", SkipTests: false}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s: table", TagPrefix)}},
		{name: "valid with spaces trimmed", expectedTable: Table{Name: "table", SkipTests: false}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s:   table  ", TagPrefix)}},
		{name: "valid with unknown flags", expectedTable: Table{Name: "table", SkipTests: false}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s: table;unkown;flags", TagPrefix)}},
		{name: "valid with skip tests flag", expectedTable: Table{Name: "table", SkipTests: true}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s: table;    skip tests  ", TagPrefix)}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			comments := make([]*ast.Comment, len(tt.comments))
			for i, c := range tt.comments {
				comments[i] = &ast.Comment{Text: c}
			}
			cg := &ast.CommentGroup{List: comments}

			tab := Table{}
			err := tab.ParseTableName(cg)

			require.Equal(t, tt.expectedErr == nil, err == nil)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedTable, tab)
			} else {
				assert.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}

func TestReconcileRelationships(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		column1 := &Column{Name: "col1"}
		table1 := &Table{Name: "table1", Columns: []*Column{column1}}

		column2 := &Column{Name: "col2", fk: "table1.col1"}
		table2 := &Table{Name: "table2", Columns: []*Column{column2}}

		dbModel := DBModel{Tables: []*Table{table1, table2}}
		err := dbModel.ReconcileRelationships()
		require.NoError(t, err)
		require.NotNil(t, column2.ForeignKey)
		assert.Equal(t, column1, column2.ForeignKey)
	})

	t.Run("invalid - FKTableAndColumn parsing", func(t *testing.T) {
		column1 := &Column{Name: "col1"}
		table1 := &Table{Name: "table1", Columns: []*Column{column1}}

		column2 := &Column{Name: "col2", fk: ".col1"}
		table2 := &Table{Name: "table2", Columns: []*Column{column2}}

		dbModel := DBModel{Tables: []*Table{table1, table2}}
		err := dbModel.ReconcileRelationships()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFKTableEmpty)
	})

	t.Run("invalid - table not in db model", func(t *testing.T) {
		column1 := &Column{Name: "col1"}
		table1 := &Table{Name: "table1", Columns: []*Column{column1}}

		column2 := &Column{Name: "col2", fk: "table3.col1"}
		table2 := &Table{Name: "table2", Columns: []*Column{column2}}

		dbModel := DBModel{Tables: []*Table{table1, table2}}
		err := dbModel.ReconcileRelationships()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFKTableNotFoundInModel)
	})

	t.Run("invalid - column not found in referenced table", func(t *testing.T) {
		column1 := &Column{Name: "col1"}
		table1 := &Table{Name: "table1", Columns: []*Column{column1}}

		column2 := &Column{Name: "col2", fk: "table1.col4"}
		table2 := &Table{Name: "table2", Columns: []*Column{column2}}

		dbModel := DBModel{Tables: []*Table{table1, table2}}
		err := dbModel.ReconcileRelationships()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrColumnNotFound)
	})
}

func TestNewDBModel_HappyPath(t *testing.T) {
	content, err := format.Source([]byte(strings.Join([]string{
		"package main",
		"import \"database/sql\"",
		"// gosqlgen: table1; skip tests",
		"type Table1 struct {",
		"Id int `gosqlgen:\"id; int; pk;ai\"`",
		"Name string `gosqlgen:\"name; varchar(255); bk\"`",
		"deleted_at sql.NullTime `gosqlgen:\"deleted_at; datetime; sd\"`",
		"}",
		"// gosqlgen: table2",
		"type Table2 struct {",
		"Id int `gosqlgen:\"id; int; pk;ai\"`",
		"Table1Id int `gosqlgen:\"table1_id; int;fk table1.id\"`",
		"}",
	}, "\n")))
	require.NoError(t, err)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	require.NoError(t, err)

	dbModel, err := NewDBModel(fset, f)
	require.NoError(t, err)

	require.NotNil(t, dbModel)
	assert.Len(t, dbModel.Tables, 2)

	// tables are sorted by name
	t1 := dbModel.Tables[0]
	t2 := dbModel.Tables[1]

	assert.Equal(t, "table1", t1.Name)
	assert.True(t, t1.SkipTests)
	assert.Len(t, t1.Columns, 3)

	columnCompare := func(same bool, typeString string, expected, tested Column) {
		compFunc := assert.Equal
		if !same {
			compFunc = assert.NotEqual
		}

		compFunc(t, expected.Name, tested.Name)
		compFunc(t, expected.FieldName, tested.FieldName)
		compFunc(t, expected.PrimaryKey, tested.PrimaryKey)
		compFunc(t, expected.ForeignKey, tested.ForeignKey)
		compFunc(t, expected.Table, tested.Table)
		compFunc(t, typeString, tested.Type.String())
		compFunc(t, expected.SoftDelete, tested.SoftDelete)
		compFunc(t, expected.BusinessKey, tested.BusinessKey)
		compFunc(t, expected.AutoIncrement, tested.AutoIncrement)
	}

	// Table: table1, Column: id
	id, err := t1.GetColumn("id")
	require.NoError(t, err)
	require.NotNil(t, id)
	columnCompare(true, "int", Column{Name: "id", FieldName: "Id", PrimaryKey: true, AutoIncrement: true, Table: t1, Type: types.Typ[types.Int], TestValuer: valuerNumeric{}}, *id)

	// Table: table1, Column: name
	name, err := t1.GetColumn("name")
	require.NoError(t, err)
	require.NotNil(t, name)
	columnCompare(true, "string", Column{Name: "name", FieldName: "Name", BusinessKey: true, Table: t1, Type: types.Typ[types.String]}, *name)

	// Table: table1, Column: name
	deletedAt, err := t1.GetColumn("deleted_at")
	require.NoError(t, err)
	require.NotNil(t, deletedAt)

	assert.Equal(t, "database/sql.NullTime", deletedAt.Type.String())
	columnCompare(true, "database/sql.NullTime", Column{Name: "deleted_at", FieldName: "deleted_at", SoftDelete: true, Table: t1}, *deletedAt)

	assert.Equal(t, "table2", t2.Name)
	assert.False(t, t2.SkipTests)

	// Table: table2, Column: id
	id2, err := t2.GetColumn("id")
	require.NoError(t, err)
	require.NotNil(t, id2)
	columnCompare(true, "int", Column{Name: "id", FieldName: "Id", PrimaryKey: true, AutoIncrement: true, Table: t2, Type: types.Typ[types.Int]}, *id2)

	table1Id, err := t2.GetColumn("table1_id")
	require.NoError(t, err)
	require.NotNil(t, table1Id)
	columnCompare(true, "int", Column{Name: "table1_id", FieldName: "Table1Id", ForeignKey: id, Table: t2, Type: types.Typ[types.Int], fk: "table1.id"}, *table1Id)
}

func TestNewDBModel_SadPath(t *testing.T) {
	cases := []struct {
		name          string
		content       string
		expectedError error
	}{
		{name: "table name parsing", content: strings.Join([]string{
			"package main",
			"// gosqlgen: ",
			"type Table1 struct {",
			"Id int `gosqlgen:\"id;pk;ai\"`",
			"}"}, "\n"),
			expectedError: ErrEmptyTablename},
		{name: "valid struct with no table annotation should be skipped", content: strings.Join([]string{
			"package main",
			"type Table1 struct {",
			"Id int `gosqlgen:\"id;pk;ai\"`",
			"}"}, "\n"),
			expectedError: nil},
		{name: "no tag found for column", content: strings.Join([]string{
			"package main",
			"// gosqlgen: table",
			"type Table1 struct {",
			"Id int",
			"}"}, "\n"),
			expectedError: ErrNoColumnTag},

		{name: "column constructor error", content: strings.Join([]string{
			"package main",
			"// gosqlgen: table",
			"type Table1 struct {",
			"Id int `gosqlgen:\"\"`",
			"}"}, "\n"),
			expectedError: ErrEmptyTag},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", tt.content, parser.ParseComments)
			require.NoError(t, err)
			_, err = NewDBModel(fset, f)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTablePkAndBk(t *testing.T) {
	cases := []struct {
		name          string
		inputTable    *Table
		expectedPK    []*Column
		expectedBK    []*Column
		expectedError error
	}{
		{
			name: "Table with Primary Key Only",
			inputTable: &Table{
				Columns: []*Column{
					{Name: "id", PrimaryKey: true},
					{Name: "name"},
				},
			},
			expectedPK: []*Column{
				{Name: "id", PrimaryKey: true},
			},
			expectedBK:    []*Column{},
			expectedError: nil,
		},
		{
			name: "Table with Primary and Business Keys",
			inputTable: &Table{
				Columns: []*Column{
					{Name: "id", PrimaryKey: true},
					{Name: "email", BusinessKey: true},
					{Name: "name"},
				},
			},
			expectedPK: []*Column{
				{Name: "id", PrimaryKey: true},
			},
			expectedBK: []*Column{
				{Name: "email", BusinessKey: true},
			},
			expectedError: nil,
		},
		{
			name: "Table with no Primary Key",
			inputTable: &Table{
				Columns: []*Column{
					{Name: "name"},
					{Name: "address"},
				},
			},
			expectedPK:    nil,
			expectedBK:    nil,
			expectedError: ErrNoPrimaryKey,
		},
		{
			name: "Table with empty columns",
			inputTable: &Table{
				Columns: []*Column{},
			},
			expectedPK:    nil,
			expectedBK:    nil,
			expectedError: ErrNoPrimaryKey,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			pk, bk, err := tt.inputTable.PkAndBk()
			assert.Equal(t, tt.expectedError == nil, err == nil)
			if err == nil {
				assert.Len(t, pk, len(tt.expectedPK))
				assert.Len(t, bk, len(tt.expectedBK))

				for i, col := range pk {
					assert.Equal(t, tt.expectedPK[i].Name, col.Name)
					assert.Equal(t, tt.expectedPK[i].PrimaryKey, col.PrimaryKey)
					assert.Equal(t, tt.expectedPK[i].BusinessKey, col.BusinessKey)
				}

				for i, col := range bk {
					assert.Equal(t, tt.expectedBK[i].Name, col.Name)
					assert.Equal(t, tt.expectedBK[i].PrimaryKey, col.PrimaryKey)
					assert.Equal(t, tt.expectedBK[i].BusinessKey, col.BusinessKey)
				}
			} else {
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}
