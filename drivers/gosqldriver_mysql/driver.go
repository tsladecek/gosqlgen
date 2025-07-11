package gosqldrivermysql

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/tsladecek/gosqlgen"
)

type driver struct {
	getTemplate *template.Template
}

const getTemplate = `
func ({{ .ObjName }} {{.StructName}}) {{.MethodName}}(ctx context.Context, db dbExecutor, {{ range .Keys }}{{.Name}} {{.Type}},{{ end }}) ({{.StructName}}, error) {
	err := db.QueryRowContext(ctx, "SELECT {{ .QueryColumns }} FROM {{ .TableName }} WHERE {{ .QueryCond }}", {{ .QueryCondValues }}).Scan({{ .ScanColumns }})
	
	if err != nil {
		return {{.StructName}}{}, err
	}

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
	queryColumns := make([]string, len(table.Columns))
	scanColumns := make([]string, len(table.Columns))
	queryCond := make([]string, 0, len(keys))
	queryCondValues := make([]string, 0, len(keys))
	objName := "t"
	for i, c := range table.Columns {
		queryColumns[i] = c.Name
		scanColumns[i] = fmt.Sprintf("&%s.%s", objName, c.FieldName)

		if c.SoftDelete {
			// TODO
		}
	}

	for _, c := range keys {
		queryCond = append(queryCond, fmt.Sprintf("%s = ?", c.Name))
		queryCondValues = append(queryCondValues, fmt.Sprintf("%s", c.Name))
	}

	data := make(map[string]any)
	data["ObjName"] = objName
	data["StructName"] = table.StructName
	data["MethodName"] = methodName
	data["Keys"] = keys
	data["QueryColumns"] = strings.Join(queryColumns, ", ")
	data["ScanColumns"] = strings.Join(scanColumns, ", ")
	data["QueryCond"] = strings.Join(queryCond, " AND ")
	data["QueryCondValues"] = strings.Join(queryCondValues, ", ")
	data["TableName"] = table.Name
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
