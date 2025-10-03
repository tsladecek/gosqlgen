package gosqlgen

import (
	"go/types"
	"io"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTestSuite struct {
	getInsert int
	update    int
	delete    int
}

func (ts *mockTestSuite) GetInsert(w io.Writer, table *Table) error {
	ts.getInsert++
	return nil
}
func (ts *mockTestSuite) Update(w io.Writer, table *Table) error {
	ts.update++
	return nil
}
func (ts *mockTestSuite) Delete(w io.Writer, table *Table) error {
	ts.delete++
	return nil
}

func (ts *mockTestSuite) ExecuteTemplate(w io.Writer, tmpl string, data any) error {
	return nil
}

type mockDriver struct {
	get    int
	create int
	update int
	delete int
}

func (d *mockDriver) Get(w io.Writer, table *Table, keys []*Column, methodName string) error {
	d.get++
	return nil
}

func (d *mockDriver) Create(w io.Writer, table *Table, methodName string) error {
	d.create++
	return nil
}

func (d *mockDriver) Update(w io.Writer, table *Table, keys []*Column, methodName string) error {
	d.update++
	return nil
}
func (d *mockDriver) Delete(w io.Writer, table *Table, keys []*Column, methodName string) error {
	d.delete++
	return nil
}

func TestCreateTemplates_MethodCalls(t *testing.T) {
	table := &Table{Columns: []*Column{{Name: "id", PrimaryKey: true, Type: types.Typ[types.String]}, {Name: "name", Type: types.Typ[types.String]}}}
	model := &DBModel{Tables: []*Table{table}}

	cases := []struct {
		name              string
		driver            *mockDriver
		testSuite         *mockTestSuite
		tableFlags        []TableFlag
		expectedDriver    mockDriver
		expectedTestSuite mockTestSuite
		expectedErr       string
	}{
		{name: "valid", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 1, delete: 1}, expectedTestSuite: mockTestSuite{getInsert: 1, update: 1, delete: 1}},
		{name: "valid - ignore", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 0, create: 0, update: 0, delete: 0}, expectedTestSuite: mockTestSuite{getInsert: 0, update: 0, delete: 0}, tableFlags: []TableFlag{TableFlagIgnore}},
		{name: "valid - ignore update", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 0, delete: 1}, expectedTestSuite: mockTestSuite{getInsert: 1, update: 0, delete: 1}, tableFlags: []TableFlag{TableFlagIgnoreUpdate}},
		{name: "valid - ignore delete", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 1, delete: 0}, expectedTestSuite: mockTestSuite{getInsert: 1, update: 1, delete: 0}, tableFlags: []TableFlag{TableFlagIgnoreDelete}},
		{name: "valid - ignore update and delete", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 0, delete: 0}, expectedTestSuite: mockTestSuite{getInsert: 1, update: 0, delete: 0}, tableFlags: []TableFlag{TableFlagIgnoreDelete, TableFlagIgnoreUpdate}},
		{name: "valid - ignore test", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 1, delete: 1}, expectedTestSuite: mockTestSuite{getInsert: 0, update: 0, delete: 0}, tableFlags: []TableFlag{TableFlagIgnoreTest}},
		{name: "valid - ignore test update", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 1, delete: 1}, expectedTestSuite: mockTestSuite{getInsert: 1, update: 0, delete: 1}, tableFlags: []TableFlag{TableFlagIgnoreTestUpdate}},
		{name: "valid - ignore test delete", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 1, delete: 1}, expectedTestSuite: mockTestSuite{getInsert: 1, update: 1, delete: 0}, tableFlags: []TableFlag{TableFlagIgnoreTestDelete}},
		{name: "valid - ignore update and test delete", driver: &mockDriver{}, testSuite: &mockTestSuite{}, expectedDriver: mockDriver{get: 1, create: 1, update: 0, delete: 1}, expectedTestSuite: mockTestSuite{getInsert: 1, update: 0, delete: 0}, tableFlags: []TableFlag{TableFlagIgnoreTestDelete, TableFlagIgnoreUpdate}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			d := t.TempDir()
			out := path.Join(d, "out.go")
			outTest := path.Join(d, "out_test.go")

			table.Flags = tt.tableFlags

			err := CreateTemplates(tt.driver, model, tt.testSuite, out, outTest)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedDriver, *tt.driver)
				assert.Equal(t, tt.expectedTestSuite, *tt.testSuite)
			}
		})
	}
}
