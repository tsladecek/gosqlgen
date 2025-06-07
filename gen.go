//go:build ignore

package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Column struct {
	Name       string
	PrimaryKey bool
	ForeignKey *Column
	Table      *Table
	Type       ast.Expr

	fk string
}

func (c *Column) FKTableAndColumn() (string, string, error) {
	m := strings.Split(c.fk, ".")
	if len(m) != 2 {
		return "", "", fmt.Errorf("Invalid FK format %s", c.fk)
	}
	return m[0], m[1], nil
}

func NewColumn(tag string) (Column, error) {
	if tag == "" {
		return Column{}, nil
	}
	items := strings.Split(tag, ",")

	c := Column{}
	c.Name = items[0]

	if len(items) == 1 {
		return c, nil
	} else if len(items) > 2 {
		return Column{}, errors.New("Invalid Column Spec: %s. Expecting at most two comma separated fields <column name>[,<id|fk table.column>]")
	}

	m := strings.TrimSpace(items[1])
	if m == "pk" {
		c.PrimaryKey = true
	} else if strings.HasPrefix(m, "fk") {
		fkFields := strings.Split(m, " ")
		if len(fkFields) != 2 {
			return Column{}, errors.New("Invalid Foreign key spec. Must be in format: fk table.column")
		}

		c.fk = fkFields[1]
	}

	return c, nil
}

type Table struct {
	Name       string
	StructName string
	Columns    []*Column
}

func (t *Table) GetColumn(columnName string) (*Column, error) {
	for _, c := range t.Columns {
		if c.Name == columnName {
			return c, nil
		}
	}

	return nil, fmt.Errorf("Column %s not found", columnName)
}

type DBModel struct {
	Tables []*Table
}

func (d *DBModel) ReconcileRelationships() error {
	tmap := make(map[string]*Table, len(d.Tables))
	for _, t := range d.Tables {
		tmap[t.Name] = t
	}

	for _, t := range d.Tables {
		for _, c := range t.Columns {
			if c.fk != "" {
				table, column, err := c.FKTableAndColumn()

				if err != nil {
					return err
				}

				tt, ok := tmap[table]
				if !ok {
					return fmt.Errorf("Table %s not found in spec", table)
				}

				col, err := tt.GetColumn(column)
				if err != nil {
					return fmt.Errorf("Column %s not found in table %s", column, table)
				}

				c.ForeignKey = col
			}
		}
	}
	return nil
}

func TableName(cgroup *ast.CommentGroup) (string, error) {
	tableName := ""
	if cgroup != nil {
		for _, c := range cgroup.List {
			if strings.HasPrefix(c.Text, "// sql: ") {
				tableName = strings.TrimPrefix(c.Text, "// sql: ")
			}
		}
	}

	return tableName, nil
}

func ExtractTagContent(tagName, input string) (string, error) {
	prefix := fmt.Sprintf(`%s:"`, tagName)
	suffix := `"`

	startIndex := strings.Index(input, prefix)
	if startIndex == -1 {
		return "", fmt.Errorf("prefix '%s' not found", prefix)
	}

	startIndex += len(prefix) // Move past the prefix

	endIndex := strings.Index(input[startIndex:], suffix)
	if endIndex == -1 {
		return "", fmt.Errorf("closing quote '%s' not found after prefix", suffix)
	}

	return input[startIndex : startIndex+endIndex], nil
}

func main() {
	filename := os.Getenv("GOFILE")
	if filename == "" {
		fmt.Println("Error: GOFILE environment variable not set.")
		os.Exit(1)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	dbModel := DBModel{Tables: make([]*Table, 0)}
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		table := Table{}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			x, isStruct := typeSpec.Type.(*ast.StructType)
			if !isStruct {
				continue
			}

			table.Name, err = TableName(genDecl.Doc)
			if err != nil {
				return
			}

			table.StructName = typeSpec.Name.Name
			if x.Fields != nil {
				for _, fff := range x.Fields.List {
					tag, err := ExtractTagContent("sql", fff.Tag.Value)
					if err != nil {
						return
					}

					column, err := NewColumn(tag)
					if err != nil {
						return
					}
					column.Table = &table
					column.Type = fff.Type
					table.Columns = append(table.Columns, &column)
				}
			}
		}

		dbModel.Tables = append(dbModel.Tables, &table)
	}

	err = dbModel.ReconcileRelationships()
	if err != nil {
		println(err.Error())
		return
	}
	for _, table := range dbModel.Tables {
		fmt.Printf("%+v\n", table)
		for _, c := range table.Columns {
			fmt.Printf("%+v\n", c)
		}
	}
}
