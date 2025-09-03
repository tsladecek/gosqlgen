package gosqlgen

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"
)

func (t *Table) testInsert(w io.Writer) {
	d := []string{}

	for _, c := range t.Columns {
		if c.ForeignKey == nil {
			continue
		}

		c.ForeignKey.Table.testInsert(w)
		d = append(d, fmt.Sprintf("%s: tbl_%s.%s", c.FieldName, c.ForeignKey.Table.Name, c.ForeignKey.FieldName))
	}

	fmt.Fprintf(w, `tbl_%s := %s{%s}
		err = tbl_%s.insert(ctx, testDb)
		require.NoError(t, err)
		`, t.Name, t.StructName, strings.Join(d, ", "), t.Name)
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
		return fmt.Errorf("Could not parse primary and business keys from table: %w", err)
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
	table.testInsert(&inserts)

	data["Inserts"] = inserts.String()

	ts.testTemplate.ExecuteTemplate(w, "main", data)
	return nil
}
