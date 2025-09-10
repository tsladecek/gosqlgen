package gosqlgen

import (
	"bytes"
	"crypto/rand"
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
	varName string
	data    []insertedValue
}

func (t *Table) testInsert(w io.Writer, previouslyInserted *insertedTable) (*insertedTable, error) {
	d := []string{}
	it := &insertedTable{varName: fmt.Sprintf("tbl_%s_%s", t.Name, rand.Text()[:8]), data: make([]insertedValue, 0)}

	for _, c := range t.Columns {
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
	}

	fmt.Fprintf(w, `%s := %s{%s}
		err = %s.insert(ctx, testDb)
		require.NoError(t, err)
		`, it.varName, t.StructName, strings.Join(d, ", "), it.varName)

	return it, nil
}

type testSuite struct {
	testTemplate *template.Template
}

//go:embed testTemplate.tmpl
var testTemplateFS embed.FS

func NewTestSuite() (testSuite, error) {
	tmpl, err := template.ParseFS(testTemplateFS, "*.tmpl")

	if err != nil {
		return testSuite{}, err
	}

	return testSuite{testTemplate: tmpl}, nil
}

func updatedValues(previouslyInserted *insertedTable) (string, *insertedTable, error) {
	res := []string{}
	newInserted := &insertedTable{varName: previouslyInserted.varName, data: make([]insertedValue, 0)}
	for _, v := range previouslyInserted.data {
		if v.column.PrimaryKey || v.column.BusinessKey || v.column.SoftDelete || v.column.ForeignKey != nil {
			continue
		}

		newValue, err := v.column.TestValuer.New(v.value)
		if err != nil {
			return "", nil, fmt.Errorf("%w: when infering new value for table=%s, column=%s for test update method", err, v.column.Table.Name, v.column.Name)
		}

		newValueFormatted, err := newValue.Format(v.column.Type)
		if err != nil {
			return "", nil, fmt.Errorf("%w: when formatting new value for table=%s, column=%s for test update method", err, v.column.Table.Name, v.column.Name)
		}

		res = append(res, fmt.Sprintf("%s.%s = %s", previouslyInserted.varName, v.column.FieldName, newValueFormatted))
		newInserted.data = append(newInserted.data, insertedValue{column: v.column, value: newValue})
	}

	return strings.Join(res, "\n"), newInserted, nil
}

func (ts testSuite) Generate(w io.Writer, driver Driver, table *Table) error {
	pk, bk, err := table.PkAndBk()
	if err != nil {
		return fmt.Errorf("%w: could not parse primary and business keys from table", err)
	}

	var inserts bytes.Buffer
	insertedData, err := table.testInsert(&inserts, nil)
	if err != nil {
		return fmt.Errorf("%w: when creating test inserts", err)
	}

	updatesPk, insertedData, err := updatedValues(insertedData)
	if err != nil {
		return fmt.Errorf("%w: when creating update template for pks", err)
	}
	updatesBk, insertedData, err := updatedValues(insertedData)
	if err != nil {
		return fmt.Errorf("%w: when creating update template for pks", err)
	}

	_ = insertedData

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

	ts.testTemplate.ExecuteTemplate(w, "main", data)
	return nil
}
