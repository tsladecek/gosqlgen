package gosqlgen

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"text/template"
)

const testTemplate = `
func TestGoSQLGen_{{.StructName}}(t *testing.T) {
	ctx := t.Context()
	var err error

	// Inserts{{ .Inserts }}

	// Get By Primary Keys
	gotByPk := {{.StructName}}{}
	err = gotByPk.{{.MethodGetByPrimaryKeys}}(ctx, testDb, {{ range .PrimaryKeys }}{{$.TableVarName}}.{{.FieldName}},{{end}})
	require.NoError(t, err)
	assert.Equal(t, {{.TableVarName}}, gotByPk)

	{{ if .BusinessKeys }}// Get By Business Keys
	gotByBk := {{.StructName}}{}
	err = gotByBk.{{.MethodGetByBusinessKeys}}(ctx, testDb, {{ range .BusinessKeys }}{{$.TableVarName}}.{{.FieldName}},{{end}})
	require.NoError(t, err)
	assert.Equal(t, {{.TableVarName}}, gotByBk)
	assert.Equal(t, gotByPk, gotByBk)
	{{ end }}
	
	{{ if and .UpdateableColumnsPK .UpdateableColumnsBK}}
	var gotAfterUpdate {{.StructName}}
	var u {{.StructName}}

	// Update By Primary Keys{{ range .UpdateableColumnsPK }}
	// {{.FieldName}}
	u = gotByPk
	u.{{ .FieldName }} = {{ .NewValue }}
	err = u.{{ $.MethodUpdateByPrimaryKeys }}(ctx, testDb)
	require.NoError(t, err)
	
	gotAfterUpdate = {{ $.StructName }}{}
	err = gotAfterUpdate.{{ $.MethodGetByPrimaryKeys }}(ctx, testDb, {{ range $.PrimaryKeys }}{{$.TableVarName}}.{{.FieldName}},{{end}} )
	require.NoError(t, err)

	assert.Equal(t, u.{{ .FieldName }}, gotAfterUpdate.{{ .FieldName }})
	{{ end }}
	{{ if .BusinessKeys }}// Update By Business Keys{{ range .UpdateableColumnsBK }}
	// {{.FieldName}}
	u = gotByBk
	u.{{ .FieldName }} = {{ .NewValue }}
	err = u.{{ $.MethodUpdateByBusinessKeys }}(ctx, testDb)
	require.NoError(t, err)
	
	gotAfterUpdate = {{ $.StructName }}{}
	err = gotAfterUpdate.{{ $.MethodGetByPrimaryKeys }}(ctx, testDb, {{ range $.PrimaryKeys }}{{$.TableVarName}}.{{.FieldName}},{{end}} )
	require.NoError(t, err)
	assert.Equal(t, u.{{ .FieldName }}, gotAfterUpdate.{{ .FieldName }})
	{{ end }}
	{{ end }}
	{{ end }}

	// Delete
	err = gotByPk.delete(ctx, testDb)
	require.NoError(t, err)
	gotAfterDelete := {{ $.StructName }}{}
	err = gotAfterDelete.{{.MethodGetByPrimaryKeys}}(ctx, testDb, {{ range .PrimaryKeys }}{{$.TableVarName}}.{{.FieldName}},{{end}})
	require.Error(t, err)
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
		require.NoError(t, err)
		`, t.Name, t.StructName, strings.Join(d, ", "), t.Name)
}

type testSuite struct {
	testTemplate *template.Template
}

func NewTestSuite() (testSuite, error) {
	tmpl, err := template.New("test").Parse(testTemplate)
	if err != nil {
		return testSuite{}, err
	}

	return testSuite{testTemplate: tmpl}, nil
}

type updatetableColumn struct {
	FieldName string
	NewValue  any
}

func newUpdateableColumn(c *Column) (updatetableColumn, error) {
	t, err := c.TypeString()
	if err != nil {
		return updatetableColumn{}, fmt.Errorf("Could not infer type of column")
	}

	switch t {
	case "int", "int8", "int16", "int32", "int64":
		return updatetableColumn{FieldName: c.FieldName, NewValue: rand.Intn(255)}, nil
	case "string":
		return updatetableColumn{FieldName: c.FieldName, NewValue: fmt.Sprintf(`"%s"`, crand.Text())}, nil
	}

	return updatetableColumn{}, fmt.Errorf("Can not infer new update value for column %s", c.Name)
}

func (ts testSuite) Generate(w io.Writer, table *Table) error {
	pk, bk, err := table.PkAndBk()
	if err != nil {
		return fmt.Errorf("Could not parse primary and business keys from table: %w", err)
	}

	updateableColumnspk := make([]updatetableColumn, 0)
	updateableColumnsbk := make([]updatetableColumn, 0)

	for _, c := range table.Columns {
		if c.PrimaryKey || c.BusinessKey || c.SoftDelete {
			continue
		}

		ucpk, err := newUpdateableColumn(c)
		if err != nil {
			return err
		}

		ucbk, err := newUpdateableColumn(c)
		if err != nil {
			return err
		}

		updateableColumnspk = append(updateableColumnspk, ucpk)
		updateableColumnsbk = append(updateableColumnsbk, ucbk)
	}

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

	ts.testTemplate.Execute(w, data)
	return nil
}
