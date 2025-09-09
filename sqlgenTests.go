package gosqlgen

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"
)

func (t *Table) testInsert(w io.Writer) error {
	d := []string{}

	for _, c := range t.Columns {
		if c.ForeignKey == nil {
			if c.SoftDelete {
				continue
			}

			v, err := c.TestValuer.New(c.TestValuer.Zero())
			if err != nil {
				return fmt.Errorf("%w: when generating new value for table=%s, column=%s", err, t.Name, c.Name)
			}
			vf, err := c.TestValuer.Format(v, c.Type.String())
			if err != nil {
				return fmt.Errorf("%w: when formating new value %t for table=%s, column=%s", err, v, t.Name, c.Name)
			}
			d = append(d, fmt.Sprintf("%s: %s", c.FieldName, vf))
			continue
		}

		c.ForeignKey.Table.testInsert(w)
		d = append(d, fmt.Sprintf("%s: tbl_%s.%s", c.FieldName, c.ForeignKey.Table.Name, c.ForeignKey.FieldName))
	}

	fmt.Fprintf(w, `tbl_%s := %s{%s}
		err = tbl_%s.insert(ctx, testDb)
		require.NoError(t, err)
		`, t.Name, t.StructName, strings.Join(d, ", "), t.Name)

	return nil
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

type updatetableColumn struct {
	FieldName string
	NewValue  any
}

func (ts testSuite) Generate(w io.Writer, driver Driver, table *Table) error {
	pk, bk, err := table.PkAndBk()
	if err != nil {
		return fmt.Errorf("%w: could not parse primary and business keys from table", err)
	}

	updateableColumnspk := make([]updatetableColumn, 0)
	updateableColumnsbk := make([]updatetableColumn, 0)

	data := make(map[string]any)
	data["StructName"] = table.StructName
	data["MethodGetByPrimaryKeys"] = MethodGetByPrimaryKeys
	data["MethodGetByBusinessKeys"] = MethodGetByBusinessKeys
	data["PrimaryKeys"] = pk
	data["BusinessKeys"] = bk
	data["TableVarName"] = fmt.Sprintf("tbl_%s", table.Name)
	data["UpdateableColumnsPK"] = updateableColumnspk
	data["UpdateableColumnsBK"] = updateableColumnsbk
	data["MethodUpdateByPrimaryKeys"] = MethodUpdateByPrimaryKeys
	data["MethodUpdateByBusinessKeys"] = MethodUpdateByBusinessKeys

	var inserts bytes.Buffer
	err = table.testInsert(&inserts)
	if err != nil {
		return fmt.Errorf("%w: when creating test inserts", err)
	}

	data["Inserts"] = inserts.String()

	ts.testTemplate.ExecuteTemplate(w, "main", data)
	return nil
}
