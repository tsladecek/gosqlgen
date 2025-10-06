package gosqlgen

import (
	"database/sql"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"
	"time"

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
		{name: "invalid - fk spec contains more than two space separated fields", tag: fmt.Sprintf(`%s:"column;pk ai;fk table col;bk;sd"`, TagPrefix), expectedErr: ErrFKSpecFieldNumber},
		{name: "invalid - fk spec contains less than two space separated fields", tag: fmt.Sprintf(`%s:"column;pk ai;fk;bk;sd"`, TagPrefix), expectedErr: ErrFKSpecFieldNumber},
		{name: "valid - column with everything", tag: fmt.Sprintf(`%s:"column;pk;ai;fk table.col;bk;sd"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true, SoftDelete: true, BusinessKey: true, AutoIncrement: true, fk: "table.col", format: stringKindBasic}},
		{name: "valid - column with everything with spaces that should be trimmed", tag: fmt.Sprintf(`%s:"   column  ;      pk;ai   ;     fk   table.col   ;  bk  ;  sd  "`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true, SoftDelete: true, BusinessKey: true, AutoIncrement: true, fk: "table.col", format: stringKindBasic}},
		{name: "valid - just pk", tag: fmt.Sprintf(`%s:"column;pk"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true, format: stringKindBasic}},
		{name: "valid - unrecognized tag is skipped", tag: fmt.Sprintf(`%s:"column;bad"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", format: stringKindBasic}},
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

func TestTableGetColumn(t *testing.T) {
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

func TestTableParseTableName(t *testing.T) {
	cases := []struct {
		name          string
		expectedErr   error
		expectedTable Table
		comments      []string
	}{
		{name: "invalid no annotation found", comments: []string{"// comment line 1", "// comment line 2"}, expectedErr: ErrNoTableTag},
		{name: "invalid empty table name", comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s:", TagPrefix)}, expectedErr: ErrEmptyTablename},
		{name: "invalid empty table name of spaces", comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s:   ", TagPrefix)}, expectedErr: ErrEmptyTablename},
		{name: "valid", expectedTable: Table{Name: "table"}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s: table", TagPrefix)}},
		{name: "valid with spaces trimmed", expectedTable: Table{Name: "table"}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s:   table  ", TagPrefix)}},
		{name: "valid with unknown flags", expectedTable: Table{Name: "table"}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s: table;unkown;flags", TagPrefix)}},
		{name: "valid with skip tests flag", expectedTable: Table{Name: "table", Flags: []TableFlag{TableFlagIgnoreDelete, TableFlagIgnoreTestUpdate}}, comments: []string{"// comment line 1", "// comment line 2", fmt.Sprintf("// %s: table;    ignore delete; ignore test update  ", TagPrefix)}},
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

func TestDBModelReconcileRelationships(t *testing.T) {
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

type customTestType struct {
	stringReturn string
}

func (t customTestType) Underlying() types.Type {
	return customTestType{stringReturn: t.stringReturn}
}
func (t customTestType) String() string {
	return t.stringReturn
}

func TestColumnInferTestValuer(t *testing.T) {
	defaultCharSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	defaultStringLength := 32
	defaultMax := float64(32)

	cases := []struct {
		name               string
		column             *Column
		expectedTestValuer TestValuer
		expectedErr        error
	}{
		{name: "string - json from type", column: &Column{Type: customTestType{stringReturn: "encoding/json.RawMessage"}}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindJSON, charSet: defaultCharSet}},
		{name: "string - json from flag", column: &Column{Type: types.Typ[types.String], format: stringKindJSON}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindJSON, charSet: defaultCharSet}},
		{name: "string - uuid", column: &Column{Type: types.Typ[types.String], format: stringKindUUID}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindUUID, charSet: defaultCharSet}},
		{name: "string - string", column: &Column{Type: types.Typ[types.String]}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: defaultCharSet}},
		{name: "string - []byte", column: &Column{Type: customTestType{stringReturn: "[]byte"}}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: defaultCharSet}},
		{name: "string - byte", column: &Column{Type: customTestType{stringReturn: "byte"}}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: defaultCharSet}},
		{name: "string - rune", column: &Column{Type: customTestType{stringReturn: "rune"}}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: defaultCharSet}},
		{name: "string - sql.NullString", column: &Column{Type: customTestType{stringReturn: "database/sql.NullString"}}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: defaultCharSet}},
		{name: "string - sql.NullByte", column: &Column{Type: customTestType{stringReturn: "database/sql.NullByte"}}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: defaultCharSet}},

		{name: "string - string, charset", column: &Column{Type: types.Typ[types.String], charSet: []rune("abcd")}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: []rune("abcd")}},
		{name: "string - string, valueset", column: &Column{Type: types.Typ[types.String], valueSet: []string{"abcd", "efgh"}}, expectedTestValuer: valuerString{length: defaultStringLength, kind: stringKindBasic, charSet: defaultCharSet, valueSet: []string{"abcd", "efgh"}}},

		{name: "numeric - integer - int", column: &Column{Type: types.Typ[types.Int]}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - int8", column: &Column{Type: types.Typ[types.Int8]}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - int16", column: &Column{Type: types.Typ[types.Int16]}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - int32", column: &Column{Type: types.Typ[types.Int32]}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - int64", column: &Column{Type: types.Typ[types.Int64]}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - sql.NullInt16", column: &Column{Type: customTestType{stringReturn: "database/sql.NullInt16"}}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - sql.NullInt32", column: &Column{Type: customTestType{stringReturn: "database/sql.NullInt32"}}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - sql.NullInt64", column: &Column{Type: customTestType{stringReturn: "database/sql.NullInt64"}}, expectedTestValuer: valuerNumeric{max: defaultMax}},
		{name: "numeric - integer - min, max", column: &Column{Type: types.Typ[types.Int], min: 6, max: 16}, expectedTestValuer: valuerNumeric{max: 16, min: 6}},
		{name: "numeric - float - float32", column: &Column{Type: types.Typ[types.Float32]}, expectedTestValuer: valuerNumeric{max: defaultMax, isFloat: true}},
		{name: "numeric - float - float64", column: &Column{Type: types.Typ[types.Float64]}, expectedTestValuer: valuerNumeric{max: defaultMax, isFloat: true}},
		{name: "numeric - float - sql.NullFloat64", column: &Column{Type: customTestType{stringReturn: "database/sql.NullFloat64"}}, expectedTestValuer: valuerNumeric{max: defaultMax, isFloat: true}},

		{name: "time - time.Time", column: &Column{Type: customTestType{stringReturn: "time.Time"}}, expectedTestValuer: valuerTime{}},
		{name: "time - sql.NullTime", column: &Column{Type: customTestType{stringReturn: "database/sql.NullTime"}}, expectedTestValuer: valuerTime{}},

		{name: "boolean - bool", column: &Column{Type: types.Typ[types.Bool]}, expectedTestValuer: valuerBoolean{}},
		{name: "boolean - sql.NullBool", column: &Column{Type: customTestType{stringReturn: "database/sql.NullBool"}}, expectedTestValuer: valuerBoolean{}},

		{name: "unsupported type", column: &Column{Type: customTestType{stringReturn: "unsupported type"}}, expectedErr: ErrUnsuportedType},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.column.inferTestValuer()
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedTestValuer, tt.column.TestValuer)
			}
		})
	}
}

func TestNewDBModel_HappyPath(t *testing.T) {
	defaultCharSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	defaultStringLength := 32
	content, err := format.Source([]byte(strings.Join([]string{
		"package main",
		"import \"database/sql\"",
		"// gosqlgen: table1; ignore test",
		"type Table1 struct {",
		"Id int `gosqlgen:\"id;pk;ai;max 16\"`",
		"Name string `gosqlgen:\"name; bk; length 8; charset (a,b,c,d)\"`",
		"deleted_at sql.NullTime `gosqlgen:\"deleted_at;sd\"`",
		"ShouldBeJSON string `gosqlgen:\"should_be_json; uuid; json;\"`", // although there is also uuid flag, format should be json as it is last
		"ShouldBeUUID string `gosqlgen:\"should_be_uuid; json; uuid\"`",  // although there is also json flag, format should be uuid as it is last
		"}",
		"// gosqlgen: table2",
		"type Table2 struct {",
		"Id int `gosqlgen:\"id;pk;ai\"`",
		"Table1Id int `gosqlgen:\"table1_id;fk table1.id\"`",
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
	assert.True(t, t1.HasFlag(TableFlagIgnoreTest))
	assert.Len(t, t1.Columns, 5)

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
		compFunc(t, expected.TestValuer, tested.TestValuer)
	}

	// Table: table1, Column: id
	id, err := t1.GetColumn("id")
	require.NoError(t, err)
	require.NotNil(t, id)
	assert.Equal(t, *id, Column{Name: "id", FieldName: "Id", PrimaryKey: true, AutoIncrement: true, Table: t1, Type: types.Typ[types.Int], max: 16, format: stringKindBasic, TestValuer: valuerNumeric{max: 16}})

	// Table: table1, Column: name
	name, err := t1.GetColumn("name")
	require.NoError(t, err)
	require.NotNil(t, name)
	assert.Equal(t, *name, Column{Name: "name", FieldName: "Name", BusinessKey: true, Table: t1, Type: types.Typ[types.String], charSet: []rune("abcd"), length: 8, format: stringKindBasic, TestValuer: valuerString{length: 8, charSet: []rune("abcd"), kind: stringKindBasic}})

	// Table: table1, Column: deleted_at
	deletedAt, err := t1.GetColumn("deleted_at")
	require.NoError(t, err)
	require.NotNil(t, deletedAt)
	assert.Equal(t, "database/sql.NullTime", deletedAt.Type.String())
	columnCompare(true, "database/sql.NullTime", Column{Name: "deleted_at", FieldName: "deleted_at", SoftDelete: true, Table: t1, TestValuer: valuerTime{}}, *deletedAt)

	// Table: table1, Column: should_be_json
	sbj, err := t1.GetColumn("should_be_json")
	require.NoError(t, err)
	require.NotNil(t, sbj)
	assert.Equal(t, *sbj, Column{Name: "should_be_json", FieldName: "ShouldBeJSON", Table: t1, TestValuer: valuerString{length: defaultStringLength, kind: stringKindJSON, charSet: defaultCharSet}, Type: types.Typ[types.String], format: stringKindJSON})

	// Table: table1, Column: should_be_uuid
	sbu, err := t1.GetColumn("should_be_uuid")
	require.NoError(t, err)
	require.NotNil(t, sbu)
	assert.Equal(t, *sbu, Column{Name: "should_be_uuid", FieldName: "ShouldBeUUID", Table: t1, TestValuer: valuerString{length: defaultStringLength, kind: stringKindUUID, charSet: defaultCharSet}, Type: types.Typ[types.String], format: stringKindUUID})

	assert.Equal(t, "table2", t2.Name)
	assert.Empty(t, t2.Flags)

	// Table: table2, Column: id
	id2, err := t2.GetColumn("id")
	require.NoError(t, err)
	require.NotNil(t, id2)
	columnCompare(true, "int", Column{Name: "id", FieldName: "Id", PrimaryKey: true, AutoIncrement: true, Table: t2, Type: types.Typ[types.Int], TestValuer: valuerNumeric{max: 32}}, *id2)

	table1Id, err := t2.GetColumn("table1_id")
	require.NoError(t, err)
	require.NotNil(t, table1Id)
	columnCompare(true, "int", Column{Name: "table1_id", FieldName: "Table1Id", ForeignKey: id, Table: t2, Type: types.Typ[types.Int], fk: "table1.id", TestValuer: valuerNumeric{max: 32}}, *table1Id)
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

func TestTagListContent(t *testing.T) {
	cases := []struct {
		name        string
		tag         string
		values      []string
		expectedErr error
	}{
		{name: "valid", tag: "tag (val1,val2)", values: []string{"val1", "val2"}},
		{name: "invalid bad position", tag: "tag tag (val1,val2) tag tag", expectedErr: ErrFlagFormat},
		{name: "valid padded", tag: "  tag   (  val1 ,  val2)  ", values: []string{"val1", "val2"}},
		{name: "valid padded deduplicated", tag: "  tag   (  val1 ,  val2)  ", values: []string{"val1", "val2"}},
		{name: "valid single char padded deduplicated", tag: "  tag   ( a , a , b )  ", values: []string{"a", "b"}},
		{name: "invalid missing start paren", tag: "tag a,b)", expectedErr: ErrFlagFormat},
		{name: "invalid missing end paren", tag: "tag (a,b", expectedErr: ErrFlagFormat},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			content, err := tagListContent(tt.tag)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.EqualValues(t, tt.values, content)
			}
		})
	}
}

func TestTagInt(t *testing.T) {
	cases := []struct {
		name        string
		tag         string
		value       int
		expectedErr error
	}{
		{name: "valid", tag: "tag 123", value: 123},
		{name: "valid padded", tag: "  tag   123", value: 123},
		{name: "invalid not an int", tag: "tag 123s", expectedErr: ErrFlagFormat},
		{name: "invalid not an int", tag: "tag 123.123", expectedErr: ErrFlagFormat},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tagInt(tt.tag)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.value, value)
			}
		})
	}
}

func TestTagFloat(t *testing.T) {
	cases := []struct {
		name        string
		tag         string
		value       float64
		expectedErr error
	}{
		{name: "valid", tag: "tag 1.23", value: 1.23},
		{name: "valid padded", tag: "  tag   1.23", value: 1.23},
		{name: "invalid not a float", tag: "tag 1.23s", expectedErr: ErrFlagFormat},
		{name: "valid int", tag: "tag 123", value: 123},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tagFloat(tt.tag)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.value, value)
			}
		})
	}

}

type mockType struct {
	typeString           string
	underlyingTypeString string
}

func (t mockType) String() string {
	return t.typeString
}

func (t mockType) Underlying() types.Type {
	return mockType{typeString: t.underlyingTypeString, underlyingTypeString: t.underlyingTypeString}
}

func TestTestValueFormat(t *testing.T) {
	cases := []struct {
		name          string
		value         TestValue
		typ           types.Type
		expectedValue string
		expectedErr   bool
	}{
		{
			name:          "Basic Int Type (int)",
			value:         TestValue{Value: 123},
			typ:           mockType{typeString: "int", underlyingTypeString: "int"},
			expectedValue: "123",
			expectedErr:   false,
		},
		{
			name:          "Aliased Int Type (int64)",
			value:         TestValue{Value: 98765},
			typ:           mockType{typeString: "MyID", underlyingTypeString: "int64"},
			expectedValue: "98765",
			expectedErr:   false,
		},
		{
			name:          "Basic Float Type (float64)",
			value:         TestValue{Value: 456.78},
			typ:           mockType{typeString: "float64", underlyingTypeString: "float64"},
			expectedValue: "456.78",
			expectedErr:   false,
		},

		// --- Numeric Types (database/sql Nullable) ---
		{
			name:          "sql.NullInt16 Type",
			value:         TestValue{Value: 16},
			typ:           mockType{typeString: "database/sql.NullInt16", underlyingTypeString: "database/sql.NullInt16"},
			expectedValue: "sql.NullInt16{Valid: true, Int16: 16}",
			expectedErr:   false,
		},
		{
			name:          "sql.NullInt64 Type (Aliased Underlying)",
			value:         TestValue{Value: 64000},
			typ:           mockType{typeString: "MyNullID", underlyingTypeString: "database/sql.NullInt64"},
			expectedValue: "sql.NullInt64{Valid: true, Int64: 64000}",
			expectedErr:   false,
		},
		{
			name:          "sql.NullFloat64 Type",
			value:         TestValue{Value: 99.125},
			typ:           mockType{typeString: "database/sql.NullFloat64", underlyingTypeString: "database/sql.NullFloat64"},
			expectedValue: "sql.NullFloat64{Valid: true, Float64: 99.125}",
			expectedErr:   false,
		},

		// --- String Types (Basic) ---
		{
			name:          "Basic String Type",
			value:         TestValue{Value: "hello world"},
			typ:           mockType{typeString: "string", underlyingTypeString: "string"},
			expectedValue: "`hello world`",
			expectedErr:   false,
		},
		{
			name:          "Byte Type ('A')",
			value:         TestValue{Value: "A"},
			typ:           mockType{typeString: "byte", underlyingTypeString: "byte"},
			expectedValue: "byte('A')",
			expectedErr:   false,
		},
		{
			name:          "Rune Type ('€') (Aliased Underlying)",
			value:         TestValue{Value: "€"},
			typ:           mockType{typeString: "CustomRune", underlyingTypeString: "rune"},
			expectedValue: "rune('€')",
			expectedErr:   false,
		},
		{
			name:          "Byte Slice Type ([]byte)",
			value:         TestValue{Value: "binary data"},
			typ:           mockType{typeString: "[]byte", underlyingTypeString: "[]byte"},
			expectedValue: "[]byte(`binary data`)",
			expectedErr:   false,
		},

		// --- String Types (database/sql Nullable) ---
		{
			name:          "sql.NullString Type",
			value:         TestValue{Value: "nullable text"},
			typ:           mockType{typeString: "database/sql.NullString", underlyingTypeString: "database/sql.NullString"},
			expectedValue: "sql.NullString{Valid: true, String: \"nullable text\"}",
			expectedErr:   false,
		},
		{
			name:          "sql.NullByte Type",
			value:         TestValue{Value: "z"},
			typ:           mockType{typeString: "database/sql.NullByte", underlyingTypeString: "database/sql.NullByte"},
			expectedValue: "sql.NullByte{Valid: true, Byte: byte('z')}",
			expectedErr:   false,
		},

		// --- Time Types ---
		{
			name:          "Time.Time Type",
			value:         TestValue{Value: time.Now()}, // Value doesn't matter, output is hardcoded
			typ:           mockType{typeString: "time.Time", underlyingTypeString: "time.Time"},
			expectedValue: "time.Now().UTC().Truncate(time.Second)",
			expectedErr:   false,
		},
		{
			name:          "sql.NullTime Type",
			value:         TestValue{Value: sql.NullTime{Valid: true, Time: time.Now()}}, // Value doesn't matter, output is hardcoded
			typ:           mockType{typeString: "database/sql.NullTime", underlyingTypeString: "database/sql.NullTime"},
			expectedValue: "sql.NullTime{Valid: true, Time: time.Now().UTC().Truncate(time.Second)}",
			expectedErr:   false,
		},

		// --- Boolean Types ---
		{
			name:          "Basic Bool Type (true)",
			value:         TestValue{Value: true},
			typ:           mockType{typeString: "bool", underlyingTypeString: "bool"},
			expectedValue: "true",
			expectedErr:   false,
		},
		{
			name:          "Basic Bool Type (false) (Underlying check)",
			value:         TestValue{Value: false},
			typ:           mockType{typeString: "CustomBool", underlyingTypeString: "bool"},
			expectedValue: "false",
			expectedErr:   false,
		},
		{
			name:          "sql.NullBool Type",
			value:         TestValue{Value: true},
			typ:           mockType{typeString: "database/sql.NullBool", underlyingTypeString: "database/sql.NullBool"},
			expectedValue: "sql.NullBool{Valid: true, Bool: true}",
			expectedErr:   false,
		},

		// --- Unsupported Type (Error Case) ---
		{
			name:          "Unsupported Type (JSON/RawMessage)",
			value:         TestValue{Value: `{"key": 1}`},
			typ:           mockType{typeString: "encoding/json.RawMessage", underlyingTypeString: "[]byte"},
			expectedValue: fmt.Sprintf("[]byte(`%s`)", `{"key": 1}`),
			expectedErr:   false,
		},
		{
			name:          "Unsupported Type (Pointers)",
			value:         TestValue{Value: nil},
			typ:           mockType{typeString: "*int", underlyingTypeString: "*int"},
			expectedValue: "",
			expectedErr:   true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := tt.value.Format(tt.typ)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedValue, formatted)
			}
		})
	}

}

func TestIsOneOfTypes(t *testing.T) {
	i := mockType{typeString: "int", underlyingTypeString: "int"}
	u := mockType{typeString: "customInt", underlyingTypeString: "int"}

	assert.True(t, IsOneOfTypes(i, []string{"int"}))
	assert.True(t, IsOneOfTypes(u, []string{"int"}))
	assert.False(t, IsOneOfTypes(u, []string{"string"}))
}

func TestIsTableFlag(t *testing.T) {
	for _, f := range []string{"ignore", "ignore update", "ignore delete", "ignore test", "ignore test update", "ignore test delete"} {
		assert.True(t, isTableFlag(f))
	}

	assert.False(t, isTableFlag("invalid"))
}

func TestTableHasFlag(t *testing.T) {
	table := Table{Flags: []TableFlag{TableFlagIgnore, TableFlagIgnoreDelete}}

	assert.True(t, table.HasFlag(TableFlagIgnore))
	assert.True(t, table.HasFlag(TableFlagIgnoreDelete))

	assert.False(t, table.HasFlag(TableFlagIgnoreUpdate))
	assert.False(t, table.HasFlag(TableFlagIgnoreTest))
	assert.False(t, table.HasFlag(TableFlagIgnoreTestUpdate))
	assert.False(t, table.HasFlag(TableFlagIgnoreTestDelete))
}

func TestTagHasPrefix(t *testing.T) {
	assert.True(t, tagHasPrefix("taag", "ta"))
	assert.True(t, tagHasPrefix("tag two", "tag"))
	assert.True(t, tagHasPrefix("  tag two", "tag"))
	assert.True(t, tagHasPrefix("  TAG two", "tag"))
	assert.False(t, tagHasPrefix("tag two", "two"))
}

func TestTagEquals(t *testing.T) {
	assert.True(t, tagEquals("tag", "tag"))
	assert.False(t, tagEquals("tag", "ttag"))
	assert.True(t, tagEquals(" tag", "tag"))
	assert.True(t, tagEquals(" tAg", "tag"))
}

func TestTagFields(t *testing.T) {
	cases := []struct {
		name           string
		tag            string
		expectedFields []string
	}{
		{name: "basic", tag: "f1 f2 f3", expectedFields: []string{"f1", "f2", "f3"}},
		{name: "additional spaces should be trimmed", tag: "  f1    f2     f3   ", expectedFields: []string{"f1", "f2", "f3"}},
		{name: "empty", tag: "", expectedFields: []string{}},
		{name: "empty with spaces", tag: "  ", expectedFields: []string{}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			f := tagFields(tt.tag)
			assert.Equal(t, tt.expectedFields, f)
		})
	}
}

func TestStringKindIsValid(t *testing.T) {
	for _, vsk := range validStringKinds {
		sk := stringKind(strings.ToLower(string(vsk)))
		valid := sk.IsValid()
		assert.True(t, valid)
		assert.Equal(t, vsk, sk)
	}

	invalid := stringKind("invalid")
	assert.False(t, invalid.IsValid())
}

func TestStringKindIsTime(t *testing.T) {
	for k, f := range timeFormats {
		sk := stringKind(k)
		tf, ok := sk.IsTime()
		require.True(t, ok)
		assert.Equal(t, f, tf)
	}
}
