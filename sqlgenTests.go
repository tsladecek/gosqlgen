package gosqlgen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

const testTemplate = `
func TestGoSQLGen_{{.StructName}}(t *testing.T) {
	ctx := t.Context()
	var err error

	{{ .Inserts }}

	// Get By Primary Keys
	gotByPk := {{.StructName}}{}
	err = gotByPk.{{.MethodGetByPrimaryKeys}}(ctx, testDb, {{ range .PrimaryKeys }}{{$.TableVarName}}.{{.FieldName}},{{end}})
	require.NotNil(t, err)
	assert.Equal(t, {{.TableVarName}}, gotByPk)

	// Get By Business Keys
	gotByBk := {{.StructName}}{}
	err = gotByBk.{{.MethodGetByBusinessKeys}}(ctx, testDb, {{ range .BusinessKeys }}{{$.TableVarName}}.{{.FieldName}},{{end}})
	require.NotNil(t, err)
	assert.Equal(t, {{.TableVarName}}, gotByBk)
	assert.Equal(t, gotByPk, gotByBk)
}
`

func (t *Table) testInsert(w io.Writer) {
	d := []string{}

	for _, c := range t.Columns {
		if c.ForeignKey == nil {
			continue
		}

		c.ForeignKey.Table.testInsert(w)
		d = append(d, fmt.Sprintf("%s: tbl_%s.%s", c.FieldName, c.ForeignKey.Table.Name, c.ForeignKey.FieldName))
	}

	fmt.Fprintf(w, `
		tbl_%s := %s{%s}
		err = tbl_%s.insert(ctx, testDb)
		require.NotNil(t, err)
		`, t.Name, t.StructName, strings.Join(d, ", "), t.Name)
}

type testSuite struct {
	getsertTestTemplate *template.Template
}

func NewTestSuite() (testSuite, error) {
	getsertTmpl, err := template.New("getsert").Parse(testTemplate)
	if err != nil {
		return testSuite{}, err
	}

	return testSuite{getsertTestTemplate: getsertTmpl}, nil
}

func (ts testSuite) Generate(w io.Writer, table *Table) error {
	pk, bk, err := table.PkAndBk()
	if err != nil {
		return fmt.Errorf("Could not parse primary and business keys from table: %w", err)
	}

	data := make(map[string]any)
	data["StructName"] = table.StructName
	data["MethodGetByPrimaryKeys"] = MethodGetByPrimaryKeys
	data["MethodGetByBusinessKeys"] = MethodGetByBusinessKeys
	data["PrimaryKeys"] = pk
	data["BusinessKeys"] = bk
	data["TableVarName"] = fmt.Sprintf("tbl_%s", table.Name)

	var inserts bytes.Buffer
	table.testInsert(&inserts)

	data["Inserts"] = inserts.String()

	ts.getsertTestTemplate.Execute(w, data)
	return nil
}
