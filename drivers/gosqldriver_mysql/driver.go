package gosqldrivermysql

import (
	"io"
	"text/template"

	"github.com/tsladecek/gosqlgen"
)

type driver struct{}

func NewDriver() driver {
	return driver{}
}

const getTemplate = `
func (t {{.StructName}}) getBy{{.ColumnName}}(ctx context.Context, db DbExecutor, id {{.ColumnType}}) ({{.StructName}}, error) {
	return nil, nil
}
`

func (d driver) Get(w io.Writer, tw io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column) error {
	tmpl, err := template.New("get").Parse(getTemplate)
	if err != nil {
		return err
	}
	data := make(map[string]string)
	data["StructName"] = table.StructName
	data["ColumnName"] = keys[0].FieldName
	data["ColumnType"] = keys[0].Type
	tmpl.Execute(w, data)
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
