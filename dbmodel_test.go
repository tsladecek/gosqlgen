package gosqlgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFKTableAndColumn(t *testing.T) {
	cases := []struct {
		name           string
		fk             string
		valid          bool
		expectedTable  string
		expectedColumn string
	}{
		{name: "valid", fk: "table.column", valid: true, expectedTable: "table", expectedColumn: "column"},
		{name: "invalid - more than one field separated by dot", fk: "table.column.column", valid: false},
		{name: "invalid - less than one field separated by dot", fk: "table", valid: false},
		{name: "invalid - empty string", fk: "", valid: false},
		{name: "invalid - table empty", fk: ".column", valid: false},
		{name: "invalid - column empty", fk: "table.", valid: false},
		{name: "invalid - table space", fk: "  .column", valid: false},
		{name: "invalid - column space", fk: "table.  ", valid: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := Column{fk: tt.fk}
			ta, co, err := c.FKTableAndColumn()
			require.Equal(t, tt.valid, err == nil)

			if tt.valid {
				assert.Equal(t, tt.expectedTable, ta)
				assert.Equal(t, tt.expectedColumn, co)
			}
		})
	}
}

func TestExtractTagContent(t *testing.T) {
	cases := []struct {
		name      string
		tagName   string
		input     string
		output    string
		shouldErr bool
	}{
		{name: "valid", tagName: "tag", input: `tag:"input"`, output: "input", shouldErr: false},
		{name: "invalid - missing tag prefix in input", tagName: "tag", input: `:"input"`, shouldErr: true},
		{name: "invalid - missing colon", tagName: "tag", input: `tag"input"`, shouldErr: true},
		{name: "invalid - double quote not exactly after colon", tagName: "tag", input: `tag:tag"input"`, shouldErr: true},
		{name: "invalid - missing start quote", tagName: "tag", input: `tag:input"`, shouldErr: true},
		{name: "invalid - missing end quote", tagName: "tag", input: `tag:"input`, shouldErr: true},
		{name: "valid - output should be space trimmed", tagName: "tag", input: `tag:"  input "`, output: "input", shouldErr: false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ExtractTagContent(tt.tagName, tt.input)
			require.Equal(t, tt.shouldErr, err != nil)

			if !tt.shouldErr {
				assert.Equal(t, tt.output, c)
			}
		})
	}
}
