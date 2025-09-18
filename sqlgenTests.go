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
	it := &insertedTable{varName: fmt.Sprintf("tbl_%s_%s", t.Name, RandomString(8, []rune("abcdefghijkl"))), data: make([]insertedValue, 0), table: t}

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
				return nil, fmt.Errorf("%w: when generating new value for table=%s, column=%s", err, t.Name, c.Name)
			}

			vf, err := v.Format(c.Type)
			if err != nil {
				return nil, fmt.Errorf("%w: when formating new value %t for table=%s, column=%s", err, v, t.Name, c.Name)
			}
			d = append(d, fmt.Sprintf("%s: %s", c.FieldName, vf))

			it.data = append(it.data, insertedValue{column: c, value: v})
			continue
		}

		insertedChild, err := c.ForeignKey.Table.testInsert(w, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: when inserting child based on FK for table=%s, column=%s", err, t.Name, c.Name)
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

	newInserted := &insertedTable{varName: previouslyInserted.varName, data: make([]insertedValue, 0)}
	for _, v := range previouslyInserted.data {
		if v.column.PrimaryKey || v.column.BusinessKey || v.column.SoftDelete {
			continue
		}

		fmt.Printf("%+v\n", v.column)

		var newValue *TestValue

		if v.column.ForeignKey != nil {
			for _, c := range previouslyInserted.children {
				if c.table != v.column.ForeignKey.Table {
					fmt.Printf("%s != %s (%p != %p)", c.table.Name, v.column.ForeignKey.Table.Name, c.table, v.column.ForeignKey.Table)
					continue
				}

				it, err := c.table.testInsert(&requiredInserts, c)
				if err != nil {
					return "", nil, fmt.Errorf("%w: when inserting table %s for testing of update on foreign key %s", err, c.table.Name, v.column.Name)
				}

				for _, iv := range it.data {
					if iv.column != v.column.ForeignKey {
						continue
					}

					newValue = &iv.value
					break
				}

				break
			}

		} else {
			nv, err := v.column.TestValuer.New(v.value)
			if err != nil {
				return "", nil, fmt.Errorf("%w: when infering new value for table=%s, column=%s for test update method", err, v.column.Table.Name, v.column.Name)
			}
			newValue = &nv
		}

		if newValue == nil {
			return "", nil, fmt.Errorf("malformed previously inserted data - new value is nil for column=%s", v.column.Name)
		}

		newValueFormatted, err := newValue.Format(v.column.Type)
		if err != nil {
			return "", nil, fmt.Errorf("%w: when formatting new value for table=%s, column=%s for test update method", err, v.column.Table.Name, v.column.Name)
		}

		res = append(res, fmt.Sprintf("%s.%s = %s", previouslyInserted.varName, v.column.FieldName, newValueFormatted))
		newInserted.data = append(newInserted.data, insertedValue{column: v.column, value: *newValue})
	}

	inserts := requiredInserts.String()
	assignments := strings.Join(res, "\n")
	return strings.Join([]string{inserts, assignments}, "\n"), newInserted, nil
}

func (ts *testSuite) newData(table *Table) (map[string]any, error) {
	pk, bk, err := table.PkAndBk()
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse primary and business keys from table", err)
	}

	var inserts bytes.Buffer
	insertedData, err := table.testInsert(&inserts, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: when creating test inserts", err)
	}

	updatesPk, insertedData, err := updatedValues(insertedData)
	if err != nil {
		return nil, fmt.Errorf("%w: when creating update template for pks", err)
	}
	updatesBk, insertedData, err := updatedValues(insertedData)
	if err != nil {
		return nil, fmt.Errorf("%w: when creating update template for pks", err)
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
		return fmt.Errorf("%w: when generating test data for get/insert method", err)
	}

	ts.templates.ExecuteTemplate(w, "getInsert", data)
	return nil
}

func (ts *testSuite) update(w io.Writer, table *Table) error {
	data, err := ts.newData(table)
	if err != nil {
		return fmt.Errorf("%w: when generating test data for update method", err)
	}
	ts.templates.ExecuteTemplate(w, "update", data)
	return nil
}

func (ts *testSuite) delete(w io.Writer, table *Table) error {
	data, err := ts.newData(table)
	if err != nil {
		return fmt.Errorf("%w: when generating test data for delete method", err)
	}

	ts.templates.ExecuteTemplate(w, "delete", data)
	return nil
}

func (ts *testSuite) Generate(w io.Writer, table *Table) error {
	var tempW bytes.Buffer

	err := ts.getInsert(&tempW, table)
	if err != nil {
		return fmt.Errorf("%w: when generating test code for get/insert methods", err)
	}

	err = ts.update(&tempW, table)
	if err != nil {
		return fmt.Errorf("%w: when generating test code for update method", err)
	}

	err = ts.delete(&tempW, table)
	if err != nil {
		return fmt.Errorf("%w: when generating test code for delete method", err)
	}

	data := make(map[string]string)
	data["Tests"] = tempW.String()
	data["StructName"] = table.StructName
	ts.templates.ExecuteTemplate(w, "main", data)

	return nil
}
