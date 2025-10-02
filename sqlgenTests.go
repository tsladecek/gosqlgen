package gosqlgen

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type insertedValue struct {
	column *Column
	value  TestValue
}

type insertedTable struct {
	varName  string
	table    *Table
	data     []insertedValue
	children []*insertedTable
}

func (t *Table) testInsert(w io.Writer, previouslyInserted *insertedTable) (*insertedTable, error) {
	d := []string{}
	it := &insertedTable{varName: fmt.Sprintf("tbl_%s_%s", t.Name, RandomString(8, []rune("abcdefghijkl"))), data: make([]insertedValue, 0), table: t, children: make([]*insertedTable, 0)}

	for _, c := range t.Columns {
		if c.PrimaryKey && c.AutoIncrement {
			continue
		}

		if c.ForeignKey == nil {
			if c.SoftDelete {
				continue
			}

			prev := c.TestValuer.Zero()

			if previouslyInserted != nil {
				for _, p := range previouslyInserted.data {
					if p.column.Name == c.Name {
						prev = p.value
					}
				}
			}

			v, err := c.TestValuer.New(prev)
			if err != nil {
				return nil, Errorf("when generating new value for table=%s, column=%s: %w", t.Name, c.Name, err)
			}

			vf, err := v.Format(c.Type)
			if err != nil {
				return nil, Errorf("when formating new value %t for table=%s, column=%s: %w", v, t.Name, c.Name, err)
			}
			d = append(d, fmt.Sprintf("%s: %s", c.FieldName, vf))

			it.data = append(it.data, insertedValue{column: c, value: v})
			continue
		}

		insertedChild, err := c.ForeignKey.Table.testInsert(w, nil)
		if err != nil {
			return nil, Errorf("when inserting child based on FK for table=%s, column=%s: %w", t.Name, c.Name, err)
		}
		d = append(d, fmt.Sprintf("%s: %s.%s", c.FieldName, insertedChild.varName, c.ForeignKey.FieldName))
		it.children = append(it.children, insertedChild)
	}

	fmt.Fprintf(w, `%s := %s{%s}
err = %s.insert(ctx, testDb)
requireNoError(t, err)
`, it.varName, t.StructName, strings.Join(d, ", "), it.varName)

	return it, nil
}

type testSuite struct {
	templates *template.Template
}

//go:embed templates
var templateFS embed.FS

func NewTestSuite() (testSuite, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.tmpl")

	if err != nil {
		return testSuite{}, err
	}

	return testSuite{templates: tmpl}, nil
}

func updatedValues(previouslyInserted *insertedTable) (string, *insertedTable, error) {
	res := []string{}
	var requiredInserts bytes.Buffer

	newInserted := &insertedTable{varName: previouslyInserted.varName, data: make([]insertedValue, 0), children: make([]*insertedTable, 0), table: previouslyInserted.table}
	for _, col := range previouslyInserted.table.Columns {
		if col.PrimaryKey && col.AutoIncrement {
			continue
		}

		if col.ForeignKey == nil {
			if col.SoftDelete {
				continue
			}

			// must be included in data
			var v *insertedValue
			for _, vv := range previouslyInserted.data {
				if vv.column == col {
					v = &vv
					break
				}
			}

			if v == nil {
				return "", nil, Errorf("object TestValue not present in previously inserted data for table=%s, column=%s", previouslyInserted.table.Name, col.Name)
			}

			newValue, err := col.TestValuer.New(v.value)
			if err != nil {
				return "", nil, Errorf("when infering new value for table=%s, column=%s for test update method: %w", col.Table.Name, col.Name, err)
			}

			newValueFormatted, err := newValue.Format(col.Type)
			if err != nil {
				return "", nil, Errorf("when formatting new value for table=%s, column=%s for test update method: %w", col.Table.Name, col.Name, err)
			}

			res = append(res, fmt.Sprintf("%s.%s = %s", previouslyInserted.varName, col.FieldName, newValueFormatted))
			newInserted.data = append(newInserted.data, insertedValue{column: col, value: newValue})
			continue
		}

		// For FK columns, the new child tables must be inserted before referencing them
		var insertedChild *insertedTable
		for _, c := range previouslyInserted.children {
			if c.table != col.ForeignKey.Table {
				continue
			}

			it, err := c.table.testInsert(&requiredInserts, c)
			if err != nil {
				return "", nil, Errorf("when inserting table %s for testing of update on foreign key %s: %w", c.table.Name, col.Name, err)
			}
			insertedChild = it
			break
		}

		if insertedChild == nil {
			return "", nil, Errorf("fk table not found in children of previously inserted table: table=%s, column=%s", col.ForeignKey.Table.Name, col.Name)
		}

		newInserted.children = append(newInserted.children, insertedChild)
		res = append(res, fmt.Sprintf("%s.%s = %s.%s", previouslyInserted.varName, col.FieldName, insertedChild.varName, col.ForeignKey.FieldName))
	}

	inserts := requiredInserts.String()
	assignments := strings.Join(res, "\n")
	return strings.TrimPrefix(strings.Join([]string{inserts, assignments}, "\n"), "\n"), newInserted, nil
}

func (ts *testSuite) newData(table *Table) (map[string]any, error) {
	pk, bk, err := table.PkAndBk()
	if err != nil {
		return nil, Errorf("could not parse primary and business keys from table: %w", err)
	}

	var inserts bytes.Buffer
	insertedData, err := table.testInsert(&inserts, nil)
	if err != nil {
		return nil, Errorf("when creating test inserts: %w", err)
	}

	updatesPk, insertedDataPK, err := updatedValues(insertedData)
	if err != nil {
		return nil, Errorf("when creating update template for pks: %w", err)
	}
	updatesBk, _, err := updatedValues(insertedDataPK)
	if err != nil {
		return nil, Errorf("when creating update template for pks: %w", err)
	}

	data := make(map[string]any)
	data["Inserts"] = inserts.String()
	data["StructName"] = table.StructName
	data["MethodGetByPrimaryKeys"] = MethodGetByPrimaryKeys
	data["MethodGetByBusinessKeys"] = MethodGetByBusinessKeys
	data["PrimaryKeys"] = pk
	data["BusinessKeys"] = bk
	data["TableVarName"] = insertedData.varName
	data["UpdatesPK"] = updatesPk
	data["UpdatesBK"] = updatesBk
	data["MethodUpdateByPrimaryKeys"] = MethodUpdateByPrimaryKeys
	data["MethodUpdateByBusinessKeys"] = MethodUpdateByBusinessKeys

	return data, nil
}

func (ts *testSuite) getInsert(w io.Writer, table *Table) error {
	data, err := ts.newData(table)
	if err != nil {
		return Errorf("when generating test data for get/insert method: %w", err)
	}

	ts.templates.ExecuteTemplate(w, "getInsert", data)
	return nil
}

func (ts *testSuite) update(w io.Writer, table *Table) error {
	data, err := ts.newData(table)
	if err != nil {
		return Errorf("when generating test data for update method: %w", err)
	}
	ts.templates.ExecuteTemplate(w, "update", data)
	return nil
}

func (ts *testSuite) delete(w io.Writer, table *Table) error {
	data, err := ts.newData(table)
	if err != nil {
		return Errorf("when generating test data for delete method: %w", err)
	}

	ts.templates.ExecuteTemplate(w, "delete", data)
	return nil
}

func (ts *testSuite) Generate(w io.Writer, table *Table) error {
	var tempW bytes.Buffer

	err := ts.getInsert(&tempW, table)
	if err != nil {
		return Errorf("when generating test code for get/insert methods: %w", err)
	}

	err = ts.update(&tempW, table)
	if err != nil {
		return Errorf("when generating test code for update method: %w", err)
	}

	err = ts.delete(&tempW, table)
	if err != nil {
		return Errorf("when generating test code for delete method: %w", err)
	}

	data := make(map[string]string)
	data["Tests"] = tempW.String()
	data["StructName"] = table.StructName
	ts.templates.ExecuteTemplate(w, "main", data)

	return nil
}
