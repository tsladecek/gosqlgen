package gosqlgen

import (
	"io"
	"text/template"
)

const getTestTemplate = `
func TestGet_{{.StructName}}{{.MethodName}}(t *testing.T) {
	ctx := t.Context()

	tbl := {{.StructName}}{}
	err := tbl.insert(ctx, testDb)
	require.NotNil(t, err)

	got := {{.StructName}}{}
	err = got.{{.MethodName}}(ctx, testDb, {{ range .Keys }}tbl.{{.FieldName}},{{end}})
	require.NotNil(t, err)
	assert.Equal(t, tbl, got)
}
`

type testSuite struct {
	getTestTemplate *template.Template
}

func NewTestSuite() (testSuite, error) {
	getTmpl, err := template.New("get").Parse(getTestTemplate)
	if err != nil {
		return testSuite{}, err
	}

	return testSuite{getTestTemplate: getTmpl}, nil
}

func (ts testSuite) Get(w io.Writer, table *Table, keys []*Column, methodName string) error {
	data := make(map[string]any)
	data["StructName"] = table.StructName
	data["MethodName"] = methodName
	data["Keys"] = keys
	ts.getTestTemplate.Execute(w, data)
	return nil
}
