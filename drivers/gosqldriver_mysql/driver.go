package gosqldrivermysql

import (
	"io"
	"text/template"

	"github.com/tsladecek/gosqlgen"
)

type driver struct {
	getTemplate *template.Template
}

const getTemplate = `
	func (t {{.StructName}}) {{.MethodName}}(ctx context.Context, db dbExecutor, {{ range .Keys }}{{.Name}} {{.Type}},{{ end }}) ({{.StructName}}, error) {
	return {{.StructName}}{}, nil
}
`

func New() (gosqlgen.Driver, error) {
	getTmpl, err := template.New("get").Parse(getTemplate)
	if err != nil {
		return driver{}, err
	}

	return driver{getTemplate: getTmpl}, nil
}

func (d driver) get(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	data := make(map[string]any)
	data["StructName"] = table.StructName
	data["MethodName"] = methodName
	data["Keys"] = keys
	d.getTemplate.Execute(w, data)
	return nil
}

func (d driver) Get(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	err := d.get(w, table, keys, methodName)
	if err != nil {
		return err
	}

	return nil
}
func (d driver) Create(w io.Writer, table *gosqlgen.Table, methodName string) error {
	return nil
}
func (d driver) Update(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	return nil
}
func (d driver) Delete(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	return nil
}

func (d driver) TestSetup(w io.Writer, dbExecutorVarName string, migrationsPath string) error {
	w.Write([]byte(`func TestNotImplemented(t *testing.T){}`))

	return nil
}
