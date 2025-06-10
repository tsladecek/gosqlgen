package gosqldrivermysql

import (
	"io"
	"text/template"

	"github.com/tsladecek/gosqlgen"
)

type driver struct {
	getTemplate     *template.Template
	getTestTemplate *template.Template
}

const getTemplate = `
func (t {{.StructName}}) getBy{{.ColumnName}}(ctx context.Context, db dbExecutor, id {{.ColumnType}}) ({{.StructName}}, error) {
	return {{.StructName}}{}, nil
}
`

const getTestTemplate = `
func Test{{.StructName}}GetBy{{.ColumnName}}(t *testing.T) {}
`

func NewDriver() (driver, error) {
	getTmpl, err := template.New("get").Parse(getTemplate)
	if err != nil {
		return driver{}, err
	}

	getTestTmpl, err := template.New("getTest").Parse(getTestTemplate)
	if err != nil {
		return driver{}, err
	}

	return driver{getTemplate: getTmpl, getTestTemplate: getTestTmpl}, nil
}

func (d driver) get(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column) error {
	data := make(map[string]string)
	data["StructName"] = table.StructName
	data["ColumnName"] = keys[0].FieldName
	data["ColumnType"] = keys[0].Type
	d.getTemplate.Execute(w, data)
	return nil
}

func (d driver) getTest(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column) error {
	data := make(map[string]string)
	data["StructName"] = table.StructName
	data["ColumnName"] = keys[0].FieldName
	d.getTestTemplate.Execute(w, data)
	return nil
}

func (d driver) Get(w io.Writer, tw io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column) error {
	err := d.get(w, table, keys)
	if err != nil {
		return err
	}

	err = d.getTest(tw, table, keys)

	if err != nil {
		return err
	}

	return nil
}
func (d driver) Create(w io.Writer, tw io.Writer, table *gosqlgen.Table) error {
	return nil
}
func (d driver) Update(w io.Writer, tw io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column) error {
	return nil
}
func (d driver) Delete(w io.Writer, tw io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column) error {
	return nil
}
