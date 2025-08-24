package gosqlgen

import (
	"fmt"
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
		{name: "invalid - tag missing sql type", tag: fmt.Sprintf(`%s:"col"`, TagPrefix), expectedErr: ErrTagFieldNumber},
		{name: "invalid - tag parsing", tag: fmt.Sprintf(`%s:col`, TagPrefix), expectedErr: ErrInvalidTagPrefix},
		{name: "invalid - fk spec contains more than two space separated fields", tag: fmt.Sprintf(`%s:"column;int;pk ai;fk table col;bk;sd"`, TagPrefix), expectedErr: ErrFKSpecFieldNumber},
		{name: "invalid - fk spec contains less than two space separated fields", tag: fmt.Sprintf(`%s:"column;int;pk ai;fk;bk;sd"`, TagPrefix), expectedErr: ErrFKSpecFieldNumber},
		{name: "valid - column with everything", tag: fmt.Sprintf(`%s:"column;int;pk ai;fk table.col;bk;sd"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true, SoftDelete: true, BusinessKey: true, AutoIncrement: true, SQLType: "int", fk: "table.col"}},
		{name: "valid - column with everything with spaces that should be trimmed", tag: fmt.Sprintf(`%s:"   column  ;  int   ;    pk ai   ;     fk table.col   ;  bk  ;  sd  "`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", PrimaryKey: true, SoftDelete: true, BusinessKey: true, AutoIncrement: true, SQLType: "int", fk: "table.col"}},
		{name: "valid - just pk", tag: fmt.Sprintf(`%s:"column;int;pk"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", SQLType: "int", PrimaryKey: true}},
		{name: "valid - unrecognized tag is skipped", tag: fmt.Sprintf(`%s:"column;int;bad"`, TagPrefix), expectedErr: nil, expectedColumn: Column{Name: "column", SQLType: "int"}},
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
