package gosqldrivermysql

import (
	crand "crypto/rand"
	"fmt"
	"go/ast"
	"io"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/tsladecek/gosqlgen"
)

type driver struct {
	getTemplate        *template.Template
	insertTemplate     *template.Template
	updateTemplate     *template.Template
	softDeleteTemplate *template.Template
	hardDeleteTemplate *template.Template
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
	_, err := db.ExecContext(ctx, "UPDATE {{.TableName}} SET {{ .ColumnPlaceholders }} WHERE {{ .KeysPlaceholders }}", {{.Values }})
	return err
}
`

const softDeleteTemplate = `
func ({{.ObjName}} *{{.StructName}}) {{.MethodName}}(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "UPDATE {{.TableName}} SET {{ .ColumnPlaceholders }} WHERE {{ .KeysPlaceholders }}", {{.Values }})
	return err
}
`

const hardDeleteTemplate = `
func ({{.ObjName}} *{{.StructName}}) {{.MethodName}}(ctx context.Context, db dbExecutor) error {
	_, err := db.ExecContext(ctx, "DELETE FROM {{.TableName}} WHERE {{ .KeysPlaceholders }}", {{ .KeysValues }})
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

	softDeleteTmpl, err := template.New("softDelete").Parse(softDeleteTemplate)
	if err != nil {
		return driver{}, err
	}

	hardDeleteTmpl, err := template.New("hardDelete").Parse(hardDeleteTemplate)
	if err != nil {
		return driver{}, err
	}

	return driver{getTemplate: getTmpl, insertTemplate: insertTmpl, updateTemplate: updateTmpl, softDeleteTemplate: softDeleteTmpl, hardDeleteTemplate: hardDeleteTmpl}, nil
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

	values := []string{}
	columnPlaceholders := []string{}

	keysPlaceholders := []string{}

	for _, col := range table.Columns {
		if col.PrimaryKey || col.BusinessKey || col.SoftDelete {
			continue
		}

		values = append(values, fmt.Sprintf("%s.%s", objName, col.FieldName))
		columnPlaceholders = append(columnPlaceholders, fmt.Sprintf("%s = ?", col.Name))
	}

	for _, k := range keys {
		values = append(values, fmt.Sprintf("%s.%s", objName, k.FieldName))
		keysPlaceholders = append(keysPlaceholders, fmt.Sprintf("%s=?", k.Name))
	}

	data["TableName"] = table.Name
	data["Values"] = strings.Join(values, ", ")
	data["ColumnPlaceholders"] = strings.Join(columnPlaceholders, ", ")

	data["KeysPlaceholders"] = strings.Join(keysPlaceholders, " AND ")

	d.updateTemplate.Execute(w, data)
	return nil
}

func (d driver) Delete(w io.Writer, table *gosqlgen.Table, keys []*gosqlgen.Column, methodName string) error {
	data := make(map[string]any)
	objName := "t"
	data["ObjName"] = objName
	data["StructName"] = table.StructName
	data["MethodName"] = methodName

	keysValues := []string{}
	keysPlaceholders := []string{}

	for _, k := range keys {
		keysValues = append(keysValues, fmt.Sprintf("%s.%s", objName, k.FieldName))
		keysPlaceholders = append(keysPlaceholders, fmt.Sprintf("%s=?", k.Name))
	}

	data["TableName"] = table.Name
	data["KeysValues"] = strings.Join(keysValues, ", ")
	data["KeysPlaceholders"] = strings.Join(keysPlaceholders, " AND ")

	softCols := make([]*gosqlgen.Column, 0)
	for _, col := range table.Columns {
		if col.SoftDelete {
			softCols = append(softCols, col)
		}
	}

	// this is hard delete
	if len(softCols) == 0 {
		d.hardDeleteTemplate.Execute(w, data)
		return nil
	}

	columnValues := []string{}
	columnPlaceholders := []string{}

	for _, col := range softCols {
		cType, err := col.TypeString()
		if err != nil {
			return fmt.Errorf("can not construct statement due to bad soft delete column type: %w", err)
		}

		switch cType {
		case "bool":
			columnValues = append(columnValues, "true")
			columnPlaceholders = append(columnPlaceholders, fmt.Sprintf("%s = ?", col.Name))
		case "sql.NullTime", "string", "time.Time":
			columnPlaceholders = append(columnPlaceholders, fmt.Sprintf("%s = CURRENT_TIMESTAMP", col.Name))
		default:
			return fmt.Errorf("Unsupported type for soft delete column %s.%s", col.Table.Name, col.Name)
		}

	}

	data["Values"] = strings.Join(slices.Concat(columnValues, keysValues), ", ")
	data["ColumnPlaceholders"] = strings.Join(columnPlaceholders, ", ")

	d.updateTemplate.Execute(w, data)
	return nil
}

func randString(maxLength int) string {
	s := crand.Text()
	return s[:min(len(s), maxLength)]
}

func (d driver) RandValue(c *gosqlgen.Column) (any, error) {
	st := c.SQLType

	if slices.Contains([]string{"bigint", "int", "int1", "int2", "int3", "int4", "int8", "integer", "smallint", "tinyint"}, strings.ToLower(st)) {
		return rand.Intn(128), nil
	}

	if slices.Contains([]string{"double", "float", "float4", "float8"}, strings.ToLower(st)) {
		return rand.Float32(), nil
	}

	if slices.Contains([]string{"bool", "boolean"}, strings.ToLower(st)) {
		return []bool{true, false}[rand.Intn(2)], nil
	}

	if slices.Contains([]string{"tinyblob", "blob", "mediumblob", "longblob", "tinytext", "text", "mediumtext", "longtext"}, strings.ToLower(st)) {
		return randString(255), nil
	}

	if after, ok := strings.CutPrefix(strings.ToLower(st), "varchar("); ok {
		if after, ok := strings.CutSuffix(after, ")"); ok {
			l, err := strconv.Atoi(after)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse varchar length to integer")
			}

			return randString(l), nil

		} else {
			return nil, fmt.Errorf("Invalid varchar sql type. Should be in format: varchar(<length>)")
		}
	}

	p := ""
	if strings.HasPrefix(strings.ToLower(st), "enum(") {
		p = st[:5]
	}

	if p != "" {
		if after, ok := strings.CutPrefix(st, p); ok {
			if after, ok := strings.CutSuffix(after, ")"); ok {
				choices := strings.Split(strings.ReplaceAll(after, "'", ""), ",")
				return strings.TrimSpace(choices[rand.Intn(len(choices))]), nil

			} else {
				return nil, fmt.Errorf("Invalid enum sql type. Should be in format: enum(<values>)")
			}
		}
	}

	return nil, fmt.Errorf("Failed to generate random value for column")
}
