package gosqlgen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

const TagPrefix = "gosqlgen"

type Column struct {
	Name          string   // name of the sql column
	FieldName     string   // name of the field in the struct
	PrimaryKey    bool     // is this a primary key column?
	ForeignKey    *Column  // address of the reference column. nil if not FK
	Table         *Table   // address of the table this Column belongs to
	Type          ast.Expr // go type of the column in the struct
	SoftDelete    bool     // does this column represent soft deletion (sd)
	BusinessKey   bool     // is this business key (bk)
	AutoIncrement bool     // is this auto incremented? Important for inserts, since this column must be fetched

	fk string
}

type Table struct {
	Name       string // name of the sql table
	StructName string // name of the struct
	Columns    []*Column
}

type DBModel struct {
	Tables      []*Table
	PackageName string
}

func (d DBModel) Debug() {
	fmt.Println("---DBModel Debug---")
	fmt.Printf("--PackageName: %s--\n", d.PackageName)
	for _, t := range d.Tables {
		fmt.Printf("--Table: Name: %s, StructName: %s--\n", t.Name, t.StructName)
		fmt.Println("Columns:")
		for _, c := range t.Columns {
			fmt.Printf("%+v\n", c)
			if c.ForeignKey != nil {
				fmt.Printf("\tFK: %+v\n", c.ForeignKey)
			}
		}

		println()
	}
}

func (c *Column) FKTableAndColumn() (string, string, error) {
	m := strings.Split(c.fk, ".")
	if len(m) != 2 {
		return "", "", fmt.Errorf("Invalid FK format %s", c.fk)
	}
	return m[0], m[1], nil
}

func ExtractTagContent(tagName, input string) (string, error) {
	prefix := fmt.Sprintf(`%s:"`, tagName)
	suffix := `"`

	startIndex := strings.Index(input, prefix)
	if startIndex == -1 {
		return "", fmt.Errorf("prefix '%s' not found", prefix)
	}

	startIndex += len(prefix)

	endIndex := strings.Index(input[startIndex:], suffix)
	if endIndex == -1 {
		return "", fmt.Errorf("closing quote '%s' not found after prefix", suffix)
	}

	return input[startIndex : startIndex+endIndex], nil
}

func NewColumn(tag string) (*Column, error) {
	tag, err := ExtractTagContent(TagPrefix, tag)

	if err != nil {
		return nil, fmt.Errorf("Invalid tag: %w", err)
	}

	if tag == "" {
		return nil, nil
	}
	items := strings.Split(tag, ",")

	c := &Column{}
	c.Name = items[0]

	if len(items) == 1 {
		return c, nil
	} else if len(items) > 2 {
		return nil, errors.New("Invalid Column Spec: %s. Expecting at most two comma separated fields <column name>[,<pk|fk table.column|bk|sd>]")
	}

	m := strings.TrimSpace(items[1])
	if m == "pk" {
		c.PrimaryKey = true
	} else if m == "pk ai" {
		c.PrimaryKey = true
		c.AutoIncrement = true
	} else if strings.HasPrefix(m, "fk") {
		fkFields := strings.Split(m, " ")
		if len(fkFields) != 2 {
			return nil, errors.New("Invalid Foreign key spec. Must be in format: fk table.column")
		}

		c.fk = fkFields[1]
	} else if m == "sd" {
		c.SoftDelete = true
	} else if m == "bk" {
		c.BusinessKey = true
	}

	return c, nil
}

func (t *Table) GetColumn(columnName string) (*Column, error) {
	for _, c := range t.Columns {
		if c.Name == columnName {
			return c, nil
		}
	}

	return nil, fmt.Errorf("Column %s not found", columnName)
}

func (t *Table) ParseTableName(cgroup *ast.CommentGroup) error {
	stripPrefix := fmt.Sprintf("// %s: ", TagPrefix)
	if cgroup != nil {
		for _, c := range cgroup.List {
			if after, ok := strings.CutPrefix(c.Text, stripPrefix); ok {
				t.Name = after
				return nil
			}
		}
	}

	return fmt.Errorf("Make sure that the struct has a doc comment line of following format: // %s:<tableName>", TagPrefix)
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

func NewDBModel(f *ast.File) (*DBModel, error) {
	dbModel := DBModel{Tables: make([]*Table, 0), PackageName: f.Name.Name}
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

			err := table.ParseTableName(genDecl.Doc)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse table name: %w", err)
			}

			table.StructName = typeSpec.Name.Name
			if x.Fields != nil {
				for _, fff := range x.Fields.List {
					column, err := NewColumn(fff.Tag.Value)
					if err != nil {
						return nil, fmt.Errorf("Failed to parse column from tag %s: %w", fff.Tag.Value, err)
					}
					column.Table = &table
					column.Type = fff.Type
					column.FieldName = fff.Names[0].Name
					table.Columns = append(table.Columns, column)
				}
			}
		}

		dbModel.Tables = append(dbModel.Tables, &table)
	}

	err := dbModel.ReconcileRelationships()
	if err != nil {
		return nil, fmt.Errorf("Failed to reconcile relationships: %w", err)
	}

	return &dbModel, nil
}

func (t *Table) PkAndBk() ([]*Column, []*Column, error) {
	pk := make([]*Column, 0)
	bk := make([]*Column, 0)

	for _, c := range t.Columns {
		if c.PrimaryKey {
			pk = append(pk, c)
		} else if c.BusinessKey {
			bk = append(bk, c)
		}
	}

	if len(pk) == 0 {
		return nil, nil, fmt.Errorf("Table %s has no primary key", t.Name)
	}

	return pk, bk, nil
}

func (c *Column) TypeString() (string, error) {
	switch t := c.Type.(type) {
	case *ast.Ident:
		return t.Name, nil
	case *ast.SelectorExpr:
		pkg, ok := t.X.(*ast.Ident)
		if !ok {
			return "", fmt.Errorf("Failed to parse type for column %s in table %s", c.Name, c.Table.Name)
		}
		return fmt.Sprintf("%s.%s", pkg.Name, t.Sel.Name), nil
	}
	return "", nil
}
