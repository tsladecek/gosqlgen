package gosqlgen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"slices"
	"strconv"
	"strings"
)

const TagPrefix = "gosqlgen"

type TestValue struct {
	Value any
}

var (
	IntegerTypes     = []string{"int", "int8", "int16", "int32", "int64"}
	IntegerTypesNull = []string{"database/sql.NullInt16", "database/sql.NullInt32", "database/sql.NullInt64"}

	FloatTypes     = []string{"float32", "float64"}
	FloatTypesNull = []string{"database/sql.NullFloat64"}

	NumericTypes           = slices.Concat(IntegerTypes, FloatTypes)
	NumericIntegerTypesAll = slices.Concat(IntegerTypes, IntegerTypesNull)
	NumericFloatTypesAll   = slices.Concat(FloatTypes, FloatTypesNull)
	NumericTypesAll        = slices.Concat(NumericIntegerTypesAll, NumericFloatTypesAll)

	StringTypes     = []string{"string", "[]byte", "byte", "rune"}
	StringTypesNull = []string{"database/sql.NullString", "database/sql.NullByte"}
	StringTypesAll  = slices.Concat(StringTypes, StringTypesNull)
	StringTypeJSON  = []string{"encoding/json.RawMessage"}

	TimeTypes     = []string{"time/Time"}
	TimeTypesNull = []string{"database/sql.NullTime"}
	TimeTypesAll  = slices.Concat(TimeTypes, TimeTypesNull)

	BooleanTypes     = []string{"bool"}
	BooleanTypesNull = []string{"database/sql.NullBool"}
	BooleanTypesAll  = slices.Concat(BooleanTypes, BooleanTypesNull)
)

func (tv TestValue) Format(columnType types.Type) (string, error) {
	t := columnType.String()
	u := columnType.Underlying().String()

	if slices.Contains(NumericTypesAll, t) || slices.Contains(NumericTypesAll, u) {
		switch {
		// Numeric
		case slices.Contains(NumericTypes, t) || slices.Contains(NumericTypes, u):
			return fmt.Sprintf("%v", tv.Value), nil
		case t == "database/sql.NullInt16" || u == "database/sql.NullInt16":
			return fmt.Sprintf("sql.NullInt16{Valid: true, Int16: %d}", tv.Value), nil
		case t == "database/sql.NullInt32" || u == "database/sql.NullInt32":
			return fmt.Sprintf("sql.NullInt32{Valid: true, Int32: %d}", tv.Value), nil
		case t == "database/sql.NullInt64" || u == "database/sql.NullInt64":
			return fmt.Sprintf("sql.NullInt16{Valid: true, Int16: %d}", tv.Value), nil
		case t == "database/sql.NullFloat64" || u == "database/sql.NullFloat64":
			return fmt.Sprintf("sql.NullFloat64{Valid: true, Float64: %d}", tv.Value), nil
		}
	} else if slices.Contains(StringTypesAll, t) || slices.Contains(StringTypesAll, u) {
		switch {
		case t == "string" || u == "string":
			return fmt.Sprintf(`"%s"`, tv.Value), nil
		case t == "byte" || u == "byte":
			return fmt.Sprintf("byte('%s')", tv.Value), nil
		case t == "rune" || u == "rune":
			return fmt.Sprintf("rune('%s')", tv.Value), nil
		case t == "[]byte" || u == "[]byte":
			return fmt.Sprintf("[]byte(`%s`)", tv.Value), nil
		case t == "sql.NullString":
			return fmt.Sprintf("sql.NullString{Valid: true, String: \"%s\"}", tv.Value), nil
		case t == "sql.NullByte":
			return fmt.Sprintf("sql.NullByte{Valid: true, Byte: byte('%s')}", tv.Value), nil
		}
	} else if slices.Contains(TimeTypesAll, t) || slices.Contains(TimeTypesAll, u) {
		switch {
		case t == "time.Time" || u == "time.Time":
			return "time.Now()", nil
		case t == "database/sql.NullTime" || u == "database/sql.NullTime":
			return "sql.NullTime{Valid: true, Time: time.Now()}", nil
		}
	} else if slices.Contains(BooleanTypesAll, t) || slices.Contains(BooleanTypesAll, u) {
		switch {
		case t == "bool" || u == "bool":
			return fmt.Sprintf("%t", tv.Value), nil
		case t == "database/sql.NullBool" || u == "database/sql.NullBool":
			return fmt.Sprintf("sql.NullBool{Valid: true, Bool: %t}", tv.Value), nil
		}
	}

	return "", fmt.Errorf("%w: unsupported type=%s (underlying=%s) for formatting", ErrValueFormat, t, u)
}

type TestValuer interface {
	New(prev TestValue) (TestValue, error)
	Zero() TestValue
}

type Column struct {
	Name          string     // name of the sql column
	FieldName     string     // name of the field in the struct
	PrimaryKey    bool       // is this a primary key column?
	ForeignKey    *Column    // address of the reference column. nil if not FK
	Table         *Table     // address of the table this Column belongs to
	Type          types.Type // go type of the column in the struct
	SoftDelete    bool       // does this column represent soft deletion (sd)
	BusinessKey   bool       // is this business key (bk)
	AutoIncrement bool       // is this auto incremented? Important for inserts, since this column must be fetched

	// useful only for test valuer
	min        float64
	max        float64
	length     int
	valueSet   []string
	charSet    []rune
	isJSON     bool
	isUUID     bool
	TestValuer TestValuer // TestValuer

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
	ErrFKFieldNumber = errors.New("expected two dot separated fields; expected format: table.column")
	ErrFKTableEmpty  = errors.New("no table specified; expected format: table.column")
	ErrFKColumnEmpty = errors.New("no column specified; expected format: table.column")

	ErrInvalidTagPrefix = errors.New("tag prefix not valid")
	ErrNoClosingQuote   = errors.New("tag not closed with quote")

	ErrEmptyTag          = errors.New("tag empty")
	ErrTagFieldNumber    = errors.New("tag must have at least one field representing column name")
	ErrFKSpecFieldNumber = errors.New("invalid Foreign key spec, must be in format: fk table.column")
	ErrFlagFieldNumber   = errors.New("invalid flag spec")
	ErrFlagFormat        = errors.New("invalid flag format")

	ErrColumnNotFound = errors.New("column not found")

	ErrEmptyTablename = errors.New("tag found in comment group but table name is empty")
	ErrNoTableTag     = errors.New("table tag not found")

	ErrFKTableNotFoundInModel = errors.New("table not found in spec when forming foreign key constraints")

	ErrNoColumnTag = errors.New("no column tag found")

	ErrNoPrimaryKey = errors.New("no primary key found for table")
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

func tagHasPrefix(tag string, prefix Flag) bool {
	return strings.HasPrefix(strings.ToLower(tag), string(prefix))
}

func tagEquals(tag string, value Flag) bool {
	return strings.EqualFold(strings.TrimSpace(tag), string(value))
}

func tagListContent(tag string) ([]string, error) {
	fields := strings.Split(tag, " ")
	if len(fields) < 2 {
		return nil, fmt.Errorf("%w: number of items in tag is less than two", ErrFlagFieldNumber)
	}

	content := strings.Join(fields[1:], " ")

	if !strings.HasPrefix(content, "(") || !strings.HasSuffix(content, ")") {
		return nil, fmt.Errorf("%w: tag content %s is not surrounded by parentheses", ErrFlagFormat, content)
	}
	res := []string{}

	for s := range strings.SplitSeq(strings.TrimSuffix(strings.TrimPrefix(fields[1], "("), ")"), ",") {
		res = append(res, strings.TrimSpace(s))
	}

	return res, nil
}

func tagInt(tag string) (int, error) {
	fields := strings.Split(tag, " ")
	if len(fields) != 2 {
		return 0, fmt.Errorf("%w: number of items in tag is not exactly two", ErrFlagFieldNumber)
	}
	n, err := strconv.Atoi(fields[1])

	if err != nil {
		return 0, err
	}
	return n, nil
}

func tagFloat(tag string) (float64, error) {
	fields := strings.Split(tag, " ")
	if len(fields) != 2 {
		return 0, fmt.Errorf("%w: number of items in tag is not exactly two", ErrFlagFieldNumber)
	}
	n, err := strconv.ParseFloat(fields[1], 64)

	if err != nil {
		return 0, err
	}
	return n, nil
}

type Flag string

const (
	FlagPrimaryKey    Flag = "pk"
	FlagBusinesKey    Flag = "bk"
	FlagSoftDelete    Flag = "sd"
	FlagForeignKey    Flag = "fk"
	FlagAutoIncrement Flag = "ai"
	FlagMin           Flag = "min"
	FlagMax           Flag = "max"
	FlagLength        Flag = "length"
	FlagJSON          Flag = "json"
	FlagUUID          Flag = "uuid"
	FlagValueSet      Flag = "valueset"
	FlagCharSet       Flag = "charset"
)

// NewColumn constructs Column from a tag. Foreign keys are stored
// in a temporary private field "fk". All relationships are reconcilled
// after all tables have been parsed
func NewColumn(tag string) (*Column, error) {
	tag, err := ExtractTagContent(TagPrefix, tag)

	if err != nil {
		return nil, fmt.Errorf("%w: tag=%s", err, tag)
	}

	if tag == "" {
		return nil, ErrEmptyTag
	}
	items := strings.Split(tag, ";")

	if len(items) < 1 {
		return nil, ErrTagFieldNumber
	}

	c := &Column{}
	c.Name = strings.TrimSpace(items[0])

	if len(items) == 1 {
		return c, nil
	}

	for _, tagItem := range items[1:] {
		m := strings.TrimSpace(tagItem)

		switch {
		case tagEquals(m, FlagAutoIncrement):
			c.AutoIncrement = true
		case tagEquals(m, FlagPrimaryKey):
			c.PrimaryKey = true
		case tagEquals(m, FlagBusinesKey):
			c.BusinessKey = true
		case tagEquals(m, FlagSoftDelete):
			c.SoftDelete = true
		case tagHasPrefix(m, FlagForeignKey):
			fkFields := strings.Split(m, " ")
			if len(fkFields) != 2 {
				return nil, ErrFKSpecFieldNumber
			}
			c.fk = fkFields[1]
		case tagEquals(m, FlagJSON):
			c.isJSON = true
		case tagEquals(m, FlagUUID):
			c.isUUID = true
		case tagHasPrefix(m, FlagMin):
			n, err := tagFloat(m)
			if err != nil {
				return nil, fmt.Errorf("%w: when parsing min, column=%s", err, c.Name)
			}
			c.min = n
		case tagHasPrefix(m, FlagMax):
			n, err := tagFloat(m)
			if err != nil {
				return nil, fmt.Errorf("%w: when parsing max, column=%s", err, c.Name)
			}
			c.max = n
		case tagHasPrefix(m, FlagLength):
			n, err := tagInt(m)
			if err != nil {
				return nil, fmt.Errorf("%w: when parsing length, column=%s", err, c.Name)
			}
			c.length = n
		case tagHasPrefix(m, FlagValueSet):
			valueSet, err := tagListContent(m)
			if err != nil {
				return nil, fmt.Errorf("%w: column=%s", err, c.Name)
			}

			c.valueSet = valueSet
		case tagHasPrefix(m, FlagCharSet):
			valueSet, err := tagListContent(m)
			if err != nil {
				return nil, fmt.Errorf("%w: column=%s", err, c.Name)
			}
			r := []rune{}

			for _, s := range valueSet {
				if len(s) != 1 {
					return nil, fmt.Errorf("%w: char must be of length 1, column=%s", ErrFlagFormat, c.Name)
				}
				r = append(r, rune(s[0]))
			}
			c.charSet = r
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

// ParseTableName expects to find a table annotation in one of struct type
// comment lines. The annotation should be in format: gosqlgen: table_name[;flags]
// The comment must be on a single line. It is expected that the code is properly
// formatted with gofmt
func (t *Table) ParseTableName(cgroup *ast.CommentGroup) error {
	stripPrefix := fmt.Sprintf("// %s:", TagPrefix)
	if cgroup != nil {
		for _, c := range cgroup.List {
			if after, ok := strings.CutPrefix(c.Text, stripPrefix); ok {
				after := strings.TrimSpace(after)
				if after == "" {
					return ErrEmptyTablename
				}

				items := strings.Split(after, ";")

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

// ReconcileRelationships loops over every parsed column in
// every table and checks if the column should be a foreign key,
// in which case it finds corresponding referenced *Column and
// stores the pointer in the ForeignKey field
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
					return fmt.Errorf("%w: fk=%s", err, c.fk)
				}

				tt, ok := tmap[table]
				if !ok {
					return fmt.Errorf("%w: table=%s", ErrFKTableNotFoundInModel, table)
				}

				col, err := tt.GetColumn(column)
				if err != nil {
					return fmt.Errorf("%w: column=%s, table=%s", err, column, table)
				}

				c.ForeignKey = col
			}
		}
	}
	return nil
}

func (c *Column) inferTestValuer() error {
	t := c.Type.String()
	u := c.Type.Underlying().String()

	switch {
	case slices.Contains(StringTypeJSON, t) || slices.Contains(StringTypeJSON, u):
		v, err := NewValuerString(c.length, stringKindJSON, c.charSet, c.valueSet)
		if err != nil {
			return err
		}

		c.TestValuer = v
		return nil

	case slices.Contains(StringTypesAll, t) || slices.Contains(StringTypesAll, u):
		kind := stringKindBasic
		if c.isJSON {
			kind = stringKindJSON
		} else if c.isUUID {
			kind = stringKindUUID
		}

		v, err := NewValuerString(c.length, kind, c.charSet, c.valueSet)
		if err != nil {
			return err
		}

		c.TestValuer = v
		return nil

	case slices.Contains(NumericIntegerTypesAll, t) || slices.Contains(NumericIntegerTypesAll, u):
		v, err := NewValuerNumeric(c.min, c.max, false)

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil

	case slices.Contains(NumericFloatTypesAll, t) || slices.Contains(NumericFloatTypesAll, u):
		v, err := NewValuerNumeric(c.min, c.max, true)

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil

	case slices.Contains(BooleanTypesAll, t) || slices.Contains(BooleanTypesAll, u):
		v, err := NewValuerBoolean()

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil

	case slices.Contains(TimeTypesAll, t) || slices.Contains(TimeTypesAll, u):
		v, err := NewValuerTime()

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil
	}

	return fmt.Errorf("unsupported type")
}

// NewDBModel parses the File and constructs the entire DBModel.
// In the first pass, all tables and columns are constructed. In the
// second, the relationships are reconcilled and finally the tables are
// sorted by their (database) name
func NewDBModel(fset *token.FileSet, f *ast.File) (*DBModel, error) {
	info := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	var conf types.Config
	conf.Importer = importer.Default()
	_, err := conf.Check("", fset, []*ast.File{f}, &info)
	if err != nil {
		return nil, err
	}

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
				fmt.Printf("Skipped struct %s, no parseable table definition found. If this is an error, please add it in the comment above the type\n", table.StructName)
				continue MainLoop
			}

			if err != nil {
				return nil, fmt.Errorf("Failed to parse table name: %w", err)
			}

			if x.Fields != nil {
				for _, fff := range x.Fields.List {
					if fff.Tag == nil {
						return nil, fmt.Errorf("%w: table=%s", ErrNoColumnTag, table.Name)
					}

					column, err := NewColumn(fff.Tag.Value)
					if err != nil {
						return nil, fmt.Errorf("%w: table=%s", err, table.Name)
					}
					column.Table = &table
					column.Type = info.TypeOf(fff.Type)
					column.FieldName = fff.Names[0].Name
					table.Columns = append(table.Columns, column)

					err = column.inferTestValuer()
					if err != nil {
						return nil, fmt.Errorf("%w: when inferring test valuer - table=%s, column=%s", err, table.StructName, column.FieldName)
					}
				}
			}
		}

		dbModel.Tables = append(dbModel.Tables, &table)
	}

	err = dbModel.ReconcileRelationships()
	if err != nil {
		return nil, fmt.Errorf("Failed to reconcile relationships: %w", err)
	}

	slices.SortFunc(dbModel.Tables, func(a, b *Table) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	return &dbModel, nil
}

// PkAndBk returns the primary key and business key columns
// of the table. Error is returned only if no primary key was found
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
		return nil, nil, ErrNoPrimaryKey
	}

	return pk, bk, nil
}
