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

func updateableValuePostProcess(v any) any {
	if vs, ok := v.(string); ok {
		// surround with quotes if string
		return fmt.Sprintf(`"%s"`, vs)
	}

	return v
}

func (ts testSuite) Generate(w io.Writer, driver Driver, table *Table) error {
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

		vpk, err := driver.RandValue(c)
		if err != nil {
			return fmt.Errorf("Failed to infer updateable value for column %v: %w", *c, err)
		}

		vbk, err := driver.RandValue(c)
		if err != nil {
			return fmt.Errorf("Failed to infer updateable value for column %v: %w", *c, err)
		}

		ucpk := updatetableColumn{FieldName: c.FieldName, NewValue: updateableValuePostProcess(vpk)}
		ucbk := updatetableColumn{FieldName: c.FieldName, NewValue: updateableValuePostProcess(vbk)}

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
