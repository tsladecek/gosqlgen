package gosqlgen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

const getsertTestTemplate = `
func TestGetSert_{{.StructName}}{{.MethodName}}(t *testing.T) {
	ctx := t.Context()
	var err error

	{{ .Inserts }}

	got := {{.StructName}}{}
	err = got.{{.MethodName}}(ctx, testDb, {{ range .Keys }}{{$.TableVarName}}.{{.FieldName}},{{end}})
	require.NotNil(t, err)
	assert.Equal(t, {{.TableVarName}}, got)
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
	getsertTmpl, err := template.New("getsert").Parse(getsertTestTemplate)
	if err != nil {
		return testSuite{}, err
	}

	return testSuite{getsertTestTemplate: getsertTmpl}, nil
}

func (ts testSuite) Get(w io.Writer, table *Table, keys []*Column, methodName string) error {
	data := make(map[string]any)
	data["StructName"] = table.StructName
	data["MethodName"] = methodName
	data["Keys"] = keys
	data["TableVarName"] = fmt.Sprintf("tbl_%s", table.Name)

	var inserts bytes.Buffer
	table.testInsert(&inserts)

	data["Inserts"] = inserts.String()

	ts.getsertTestTemplate.Execute(w, data)
	return nil
}
