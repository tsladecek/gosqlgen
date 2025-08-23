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
