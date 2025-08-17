package gosqldrivermysql

import (
	"fmt"
	"go/ast"
	"io"
	"slices"
	"strings"
	"text/template"

	"github.com/tsladecek/gosqlgen"
)

type driver struct {
	getTemplate    *template.Template
	insertTemplate *template.Template
	updateTemplate *template.Template
	deleteTemplate *template.Template
}

const getTemplate = `
func ({{.ObjName}} *{{.StructName}}) {{.MethodName}}(ctx context.Context, db dbExecutor, {{ range .Keys }}{{.Name}} {{.Type}},{{ end }}) error {
	err := db.QueryRowContext(ctx, "SELECT {{ .QueryColumns }} FROM {{ .TableName }} WHERE {{ .QueryCond }}", {{ .QueryCondValues }}).Scan({{ .ScanColumns }})
	
	if err != nil {
		return err
	}

	return nil
}
`

const insertTemplate = `
func ({{.ObjName}} *{{.StructName}}) {{.MethodName}}(ctx context.Context, db dbExecutor) error {
	res, err := db.ExecContext(ctx, "INSERT INTO {{.TableName}} ({{.ColumnNames}}) VALUES ({{.ColumnValuesPlaceholders}})", {{.ColumnValues}})
	if err != nil {
		return err
	}

	{{if .AutoIncrementColumn }}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		t.{{.AutoIncrementColumn.FieldName}} = {{.AIColumnType}}(id)
	{{else}}
	_ = res
	{{end}}

	return nil
}
`

const updateTemplate = `
func ({{.ObjName}} *{{.StructName}}) {{.MethodName}}(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE {{.TableName}} SET {{ .ColumnPlaceholders }} WHERE {{ .KeysPlaceholders }}", {{.ColumnValues}}, {{ .KeysValues }})
	return err
}
`

func New() (gosqlgen.Driver, error) {
	getTmpl, err := template.New("get").Parse(getTemplate)
	if err != nil {
		return driver{}, err
	}

	insertTmpl, err := template.New("insert").Parse(insertTemplate)
	if err != nil {
		return driver{}, err
	}

	updateTmpl, err := template.New("update").Parse(updateTemplate)
	if err != nil {
		return driver{}, err
	}

	return driver{getTemplate: getTmpl, insertTemplate: insertTmpl, updateTemplate: updateTmpl}, nil
}

func (d driver) Get(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	queryColumns := make([]string, len(table.Columns))
	scanColumns := make([]string, len(table.Columns))
	queryCond := make([]string, 0, len(keys))
	queryCondValues := make([]string, 0, len(keys))
	objName := "t"

	for _, c := range keys {
		queryCond = append(queryCond, fmt.Sprintf("%s = ?", c.Name))
		queryCondValues = append(queryCondValues, fmt.Sprintf("%s", c.Name))
	}

	for i, c := range table.Columns {
		queryColumns[i] = c.Name
		scanColumns[i] = fmt.Sprintf("&%s.%s", objName, c.FieldName)

		if c.SoftDelete {
			cType, err := c.TypeString()
			if err != nil {
				return fmt.Errorf("can not construct statement due to bad soft delete column type: %w", err)
			}

			switch cType {
			case "bool":
				queryCond = append(queryCond, fmt.Sprintf("%s = true", c.Name))
			case "sql.NullTime":
				queryCond = append(queryCond, fmt.Sprintf("%s IS NOT NULL", c.Name))
			case "string":
				queryCond = append(queryCond, fmt.Sprintf(`%s = ?`, c.Name))
				queryCondValues = append(queryCondValues, `""`)
			}
		}
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

func (d driver) Create(w io.Writer, table *gosqlgen.Table, methodName string) error {
	data := make(map[string]any)
	objName := "t"
	data["ObjName"] = objName
	data["StructName"] = table.StructName
	data["MethodName"] = methodName

	columnNames := []string{}
	columnValues := []string{}
	columnPlaceholders := []string{}
	var aiCol *gosqlgen.Column

	for _, col := range table.Columns {
		if col.PrimaryKey {
			if col.AutoIncrement {
				aiCol = col
			}
			continue
		}

		if col.SoftDelete {
			continue
		}

		columnNames = append(columnNames, col.Name)
		columnValues = append(columnValues, fmt.Sprintf("%s.%s", objName, col.FieldName))
		columnPlaceholders = append(columnPlaceholders, "?")
	}

	data["TableName"] = table.Name
	data["ColumnNames"] = strings.Join(columnNames, ", ")
	data["ColumnValues"] = strings.Join(columnValues, ", ")
	data["ColumnValuesPlaceholders"] = strings.Join(columnPlaceholders, ", ")
	data["AutoIncrementColumn"] = aiCol
	data["AIColumnType"] = "int"

	if aiCol != nil {
		ai, ok := aiCol.Type.(*ast.Ident)
		if !ok {
			return fmt.Errorf("Autoincrement column %s not a basic type. Must be one of following types: int, int16, int32, int64", aiCol.Name)
		}

		if !slices.Contains([]string{"int", "int16", "int32", "int64"}, ai.Name) {
			return fmt.Errorf("Autoincrement column %s must be one of following types: int, int16, int32, int64", aiCol.Name)
		}

		data["AIColumnType"] = ai.Name
	}

	d.insertTemplate.Execute(w, data)
	return nil
}
func (d driver) Update(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	data := make(map[string]any)
	objName := "t"
	data["ObjName"] = objName
	data["StructName"] = table.StructName
	data["MethodName"] = methodName

	columnValues := []string{}
	columnPlaceholders := []string{}

	keysValues := []string{}
	keysPlaceholders := []string{}

	for _, k := range keys {
		keysValues = append(keysValues, fmt.Sprintf("%s.%s", objName, k.FieldName))
		keysPlaceholders = append(keysPlaceholders, fmt.Sprintf("%s=?", k.Name))
	}

	for _, col := range table.Columns {
		if col.PrimaryKey || col.BusinessKey || col.SoftDelete {
			continue
		}

		columnValues = append(columnValues, fmt.Sprintf("%s.%s", objName, col.FieldName))
		columnPlaceholders = append(columnPlaceholders, fmt.Sprintf("%s = ?", col.Name))
	}

	data["TableName"] = table.Name
	data["ColumnValues"] = strings.Join(columnValues, ", ")
	data["ColumnPlaceholders"] = strings.Join(columnPlaceholders, ", ")

	data["KeysValues"] = strings.Join(keysValues, ", ")
	data["KeysPlaceholders"] = strings.Join(keysPlaceholders, " AND ")

	d.updateTemplate.Execute(w, data)
	return nil
}
func (d driver) Delete(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	return nil
}
