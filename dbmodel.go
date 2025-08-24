package gosqlgen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"slices"
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
	SQLType       string   // this can be driver specific

	fk string
}

type Table struct {
	Name       string // name of the sql table
	StructName string // name of the struct
	Columns    []*Column
	SkipTests  bool
}

type DBModel struct {
	Tables      []*Table
	PackageName string
}

var (
	ErrNoTableTag = errors.New("table tag not found")

	ErrFKFieldNumber = errors.New("expected two dot separated fields; expected format: table.column")
	ErrFKTableEmpty  = errors.New("no table specified; expected format: table.column")
	ErrFKColumnEmpty = errors.New("no column specified; expected format: table.column")

	ErrInvalidTagPrefix = errors.New("tag prefix not valid")
	ErrNoClosingQuote   = errors.New("tag not closed with quote")

	ErrEmptyTag          = errors.New("tag empty")
	ErrTagFieldNumber    = errors.New("tag must have two required fields: column name and sql type (e.g. name,varchar(31))")
	ErrFKSpecFieldNumber = errors.New("invalid Foreign key spec, must be in format: fk table.column")

	ErrColumnNotFound = errors.New("column not found")
)

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

// FKTableAndColumn parses the table name and the column name
// from the fk tag specification in format "table.column"
// Returns error if there are not exactly two dot separated fields separated
func (c *Column) FKTableAndColumn() (string, string, error) {
	m := strings.Split(c.fk, ".")
	if len(m) != 2 {
		return "", "", ErrFKFieldNumber
	}

	table := strings.TrimSpace(m[0])
	column := strings.TrimSpace(m[1])

	if table == "" {
		return "", "", ErrFKTableEmpty
	}

	if column == "" {
		return "", "", ErrFKColumnEmpty
	}

	return table, column, nil
}

// ExtractTagContent extracts the content of a given tagName enclosed
// within double quotes
func ExtractTagContent(tagName, input string) (string, error) {
	prefix := fmt.Sprintf(`%s:"`, tagName)
	suffix := `"`

	startIndex := strings.Index(input, prefix)
	if startIndex == -1 {
		return "", ErrInvalidTagPrefix
	}

	startIndex += len(prefix)

	endIndex := strings.Index(input[startIndex:], suffix)
	if endIndex == -1 {
		return "", ErrNoClosingQuote
	}

	return strings.TrimSpace(input[startIndex : startIndex+endIndex]), nil
}

// NewColumn constructs Column from a tag. Foreign keys are stored
// in a temporary private field "fk". All relationships are reconcilled
// after all tables have been parsed
func NewColumn(tag string) (*Column, error) {
	tag, err := ExtractTagContent(TagPrefix, tag)

	if err != nil {
		return nil, fmt.Errorf("invalid tag: %w", err)
	}

	if tag == "" {
		return nil, ErrEmptyTag
	}
	items := strings.Split(tag, ";")

	if len(items) < 2 {
		return nil, ErrTagFieldNumber
	}

	c := &Column{}
	c.Name = strings.TrimSpace(items[0])
	c.SQLType = strings.TrimSpace(items[1])

	if len(items) == 2 {
		return c, nil
	}

	for _, tagItem := range items[2:] {
		m := strings.TrimSpace(tagItem)
		if m == "pk" {
			c.PrimaryKey = true
		} else if m == "pk ai" {
			c.PrimaryKey = true
			c.AutoIncrement = true
		} else if strings.HasPrefix(m, "fk") {
			fkFields := strings.Split(m, " ")
			if len(fkFields) != 2 {
				return nil, ErrFKSpecFieldNumber
			}
			c.fk = fkFields[1]
		} else if m == "sd" {
			c.SoftDelete = true
		} else if m == "bk" {
			c.BusinessKey = true
		}
	}
	return c, nil
}

// GetColumn loops over columns in the table and returns
// the one with matching column name. In case that no is found,
// an error is returned
func (t *Table) GetColumn(columnName string) (*Column, error) {
	for _, c := range t.Columns {
		if c.Name == columnName {
			return c, nil
		}
	}

	return nil, ErrColumnNotFound
}

func (t *Table) ParseTableName(cgroup *ast.CommentGroup) error {
	stripPrefix := fmt.Sprintf("// %s: ", TagPrefix)
	if cgroup != nil {
		for _, c := range cgroup.List {
			if after, ok := strings.CutPrefix(c.Text, stripPrefix); ok {
				items := strings.Split(after, ";")
				if len(items) == 0 {
					return fmt.Errorf("Table name must not be empty")
				}

				t.Name = strings.TrimSpace(items[0])

				if len(items) > 1 {
					for _, item := range items {
						item = strings.TrimSpace(item)
						switch item {
						case "skip tests":
							t.SkipTests = true
						}
					}
				}
				return nil
			}
		}
	}

	return ErrNoTableTag
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
					return fmt.Errorf("invalid foreign key format %s: %w", c.fk, err)
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
MainLoop:
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		table := Table{}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue MainLoop
			}

			x, isStruct := typeSpec.Type.(*ast.StructType)
			if !isStruct {
				continue MainLoop
			}
			table.StructName = typeSpec.Name.Name

			err := table.ParseTableName(genDecl.Doc)
			if errors.Is(err, ErrNoTableTag) {
				fmt.Printf("Skipped struct %s, no parseable table definition found. If this is an error, please add it in the comment above the type", table.StructName)
				continue MainLoop
			}

			if err != nil {
				return nil, fmt.Errorf("Failed to parse table name: %w", err)
			}

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

	slices.SortFunc(dbModel.Tables, func(a, b *Table) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

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
		return nil, nil, fmt.Errorf("Table %s (%s) has no primary key", t.Name, t.StructName)
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
