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
	"time"
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

	TimeTypes     = []string{"time.Time"}
	TimeTypesNull = []string{"database/sql.NullTime"}
	TimeTypesAll  = slices.Concat(TimeTypes, TimeTypesNull)

	BooleanTypes     = []string{"bool"}
	BooleanTypesNull = []string{"database/sql.NullBool"}
	BooleanTypesAll  = slices.Concat(BooleanTypes, BooleanTypesNull)
)

// IsOneOfTypes checks if given type or its underlying type is one
// of provided type strings
func IsOneOfTypes(typ types.Type, options []string) bool {
	// Each type T has an underlying type: If T is one of the predeclared boolean, numeric, or string types, or a type literal, the corresponding underlying type is T itself. Otherwise, T's underlying type is the underlying type of the type to which T refers in its declaration.
	t := typ.String()
	u := typ.Underlying().String()

	return slices.Contains(options, t) || slices.Contains(options, u)
}

func (tv TestValue) Format(columnType types.Type) (string, error) {
	t := columnType.String()
	u := columnType.Underlying().String()

	if IsOneOfTypes(columnType, NumericTypesAll) {
		switch {
		// Numeric
		case IsOneOfTypes(columnType, NumericTypes):
			return fmt.Sprintf("%v", tv.Value), nil
		case t == "database/sql.NullInt16" || u == "database/sql.NullInt16":
			return fmt.Sprintf("sql.NullInt16{Valid: true, Int16: %d}", tv.Value), nil
		case t == "database/sql.NullInt32" || u == "database/sql.NullInt32":
			return fmt.Sprintf("sql.NullInt32{Valid: true, Int32: %d}", tv.Value), nil
		case t == "database/sql.NullInt64" || u == "database/sql.NullInt64":
			return fmt.Sprintf("sql.NullInt64{Valid: true, Int64: %d}", tv.Value), nil
		case t == "database/sql.NullFloat64" || u == "database/sql.NullFloat64":
			return fmt.Sprintf("sql.NullFloat64{Valid: true, Float64: %v}", tv.Value), nil
		}
	} else if IsOneOfTypes(columnType, StringTypesAll) {
		switch {
		case t == "string" || u == "string":
			return fmt.Sprintf("`%s`", tv.Value), nil
		case t == "byte" || u == "byte":
			return fmt.Sprintf("byte('%s')", tv.Value), nil
		case t == "rune" || u == "rune":
			return fmt.Sprintf("rune('%s')", tv.Value), nil
		case t == "[]byte" || u == "[]byte":
			return fmt.Sprintf("[]byte(`%s`)", tv.Value), nil
		case t == "database/sql.NullString":
			return fmt.Sprintf("sql.NullString{Valid: true, String: \"%s\"}", tv.Value), nil
		case t == "database/sql.NullByte":
			return fmt.Sprintf("sql.NullByte{Valid: true, Byte: byte('%s')}", tv.Value), nil
		}
	} else if IsOneOfTypes(columnType, TimeTypesAll) {
		switch {
		case t == "time.Time" || u == "time.Time":
			return "time.Now().UTC().Truncate(time.Second)", nil
		case t == "database/sql.NullTime" || u == "database/sql.NullTime":
			return "sql.NullTime{Valid: true, Time: time.Now().UTC().Truncate(time.Second)}", nil
		}
	} else if IsOneOfTypes(columnType, BooleanTypesAll) {
		switch {
		case t == "bool" || u == "bool":
			return fmt.Sprintf("%t", tv.Value), nil
		case t == "database/sql.NullBool" || u == "database/sql.NullBool":
			return fmt.Sprintf("sql.NullBool{Valid: true, Bool: %t}", tv.Value), nil
		}
	}

	return "", fmt.Errorf("unsupported type=%s (underlying=%s) for formatting: %w", t, u, ErrValueFormat)
}

type TestValuer interface {
	New(prev TestValue) (TestValue, error)
	Zero() TestValue
}

type stringKind string

const (
	stringKindBasic stringKind = "basic"
	stringKindEnum  stringKind = "enum"
	stringKindJSON  stringKind = "json"
	stringKindUUID  stringKind = "uuid"
	stringKindIPV4  stringKind = "ipv4"
	stringKindIPV6  stringKind = "ipv6"

	// Time
	stringKindTimeLayout      stringKind = "Layout"
	stringKindTimeANSIC       stringKind = "ANSIC"
	stringKindTimeUnixDate    stringKind = "UnixDate"
	stringKindTimeRubyDate    stringKind = "RubyDate"
	stringKindTimeRFC822      stringKind = "RFC822"
	stringKindTimeRFC822Z     stringKind = "RFC822Z"
	stringKindTimeRFC850      stringKind = "RFC850"
	stringKindTimeRFC1123     stringKind = "RFC1123"
	stringKindTimeRFC1123Z    stringKind = "RFC1123Z"
	stringKindTimeRFC3339     stringKind = "RFC3339"
	stringKindTimeRFC3339Nano stringKind = "RFC3339Nano"
	stringKindTimeKitchen     stringKind = "Kitchen"
	stringKindTimeStamp       stringKind = "Stamp"
	stringKindTimeStampMilli  stringKind = "StampMilli"
	stringKindTimeStampMicro  stringKind = "StampMicro"
	stringKindTimeStampNano   stringKind = "StampNano"
	stringKindTimeDateTime    stringKind = "DateTime"
	stringKindTimeDateOnly    stringKind = "DateOnly"
	stringKindTimeTimeOnly    stringKind = "TimeOnly"
)

var validStringKinds = []stringKind{stringKindBasic, stringKindEnum, stringKindJSON, stringKindUUID, stringKindIPV4, stringKindIPV6, stringKindTimeLayout, stringKindTimeANSIC, stringKindTimeUnixDate, stringKindTimeRubyDate, stringKindTimeRFC822, stringKindTimeRFC822Z, stringKindTimeRFC850, stringKindTimeRFC1123, stringKindTimeRFC1123Z, stringKindTimeRFC3339, stringKindTimeRFC3339Nano, stringKindTimeKitchen, stringKindTimeStamp, stringKindTimeStampMilli, stringKindTimeStampMicro, stringKindTimeStampNano, stringKindTimeDateTime, stringKindTimeDateOnly, stringKindTimeTimeOnly}

func (sk *stringKind) IsValid() bool {
	for _, vsk := range validStringKinds {
		if strings.EqualFold(string(vsk), string(*sk)) {
			*sk = vsk
			return true
		}

	}
	return false
}

func (sk stringKind) IsTime() (string, bool) {
	tf, ok := timeFormats[string(sk)]
	return tf, ok
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
	charSet    []rune
	format     stringKind
	valueSet   []string
	TestValuer TestValuer // TestValuer

	fk string
}

type TableFlag string

const (
	TableFlagIgnore           = "ignore"
	TableFlagIgnoreUpdate     = "ignore update"
	TableFlagIgnoreDelete     = "ignore delete"
	TableFlagIgnoreTest       = "ignore test"
	TableFlagIgnoreTestUpdate = "ignore test update"
	TableFlagIgnoreTestDelete = "ignore test delete"
)

func isTableFlag(flag string) bool {
	return slices.Contains([]string{TableFlagIgnore, TableFlagIgnoreUpdate, TableFlagIgnoreDelete, TableFlagIgnoreTest, TableFlagIgnoreTestUpdate, TableFlagIgnoreTestDelete}, flag)
}

type Table struct {
	Name       string // name of the sql table
	StructName string // name of the struct
	Columns    []*Column
	Flags      []TableFlag
}

func (t *Table) HasFlag(flag TableFlag) bool {
	return slices.Contains(t.Flags, flag)
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
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(tag)), string(prefix))
}

func tagEquals(tag string, value Flag) bool {
	return strings.EqualFold(strings.TrimSpace(tag), string(value))
}

// tagListContent returns list of space trimmed fields
func tagFields(tag string) []string {
	fields := []string{}
	for c := range strings.SplitSeq(strings.TrimSpace(tag), " ") {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}

		fields = append(fields, c)
	}

	return fields
}

// tagListContent extracts the content inside parentheses and
// returns it as a slice of strings. It is space agnostic. The function
// does not check the position of the content inside the tag
func tagListContent(tag string) ([]string, error) {
	fields := tagFields(tag)

	if strings.HasPrefix(fields[0], "(") {
		return nil, Errorf("tag can not start with parenthesis: %w", ErrFlagFormat)
	}

	tagContent := strings.Join(fields[1:], "")

	if !strings.HasPrefix(tagContent, "(") || !strings.HasSuffix(tagContent, ")") {
		return nil, Errorf("tag must be surrounded with parenthesis: %w", ErrFlagFormat)
	}

	content := []string{}
	for s := range strings.SplitSeq(strings.TrimPrefix(strings.TrimSuffix(tagContent, ")"), "("), ",") {
		if slices.Contains(content, s) {
			continue
		}

		content = append(content, s)
	}

	return content, nil
}

// tagInt extracts integer from the second position in the space delimited tag fields
func tagInt(tag string) (int, error) {
	fields := tagFields(tag)
	if len(fields) != 2 {
		return 0, Errorf("number of items in tag is not exactly two: %w", ErrFlagFieldNumber)
	}
	n, err := strconv.Atoi(fields[1])

	if err != nil {
		return 0, Errorf("failed to convert to integer: %w", ErrFlagFormat)
	}
	return n, nil
}

// tagFloat extracts float from the second position in the space delimited tag fields
func tagFloat(tag string) (float64, error) {
	fields := tagFields(tag)
	if len(fields) != 2 {
		return 0, Errorf("number of items in tag is not exactly two: %w", ErrFlagFieldNumber)
	}
	n, err := strconv.ParseFloat(fields[1], 64)

	if err != nil {
		return 0, Errorf("failed to convert to float: %w", ErrFlagFormat)
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
	FlagEnum          Flag = "enum"
	FlagCharSet       Flag = "charset"
	FlagTime          Flag = "time"
	FlagIPV4          Flag = "ipv4"
	FlagIPV6          Flag = "ipv6"
)

var timeFormats = map[string]string{"Layout": time.Layout, "ANSIC": time.ANSIC, "UnixDate": time.UnixDate, "RubyDate": time.RubyDate, "RFC822": time.RFC822, "RFC822Z": time.RFC822Z, "RFC850": time.RFC850, "RFC1123": time.RFC1123, "RFC1123Z": time.RFC1123Z, "RFC3339": time.RFC3339, "RFC3339Nano": time.RFC3339Nano, "Kitchen": time.Kitchen, "Stamp": time.Stamp, "StampMilli": time.StampMilli, "StampMicro": time.StampMicro, "StampNano": time.StampNano, "DateTime": time.DateTime, "DateOnly": time.DateOnly, "TimeOnly": time.TimeOnly}

// NewColumn constructs Column from a tag. Foreign keys are stored
// in a temporary private field "fk". All relationships are reconcilled
// after all tables have been parsed
func NewColumn(tag string) (*Column, error) {
	tag, err := ExtractTagContent(TagPrefix, tag)

	if err != nil {
		return nil, Errorf("tag=%s: %w", tag, err)
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
	c.format = stringKindBasic

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
			fkFields := tagFields(m)
			if len(fkFields) != 2 {
				return nil, ErrFKSpecFieldNumber
			}
			c.fk = fkFields[1]
		case tagEquals(m, FlagJSON):
			c.format = stringKindJSON
		case tagEquals(m, FlagUUID):
			c.format = stringKindUUID
		case tagEquals(m, FlagIPV4):
			c.format = stringKindIPV4
		case tagEquals(m, FlagIPV6):
			c.format = stringKindIPV6
		case tagHasPrefix(m, FlagMin):
			n, err := tagFloat(m)
			if err != nil {
				return nil, Errorf("when parsing min, column=%s: %w", c.Name, err)
			}
			c.min = n
		case tagHasPrefix(m, FlagMax):
			n, err := tagFloat(m)
			if err != nil {
				return nil, Errorf("when parsing max, column=%s: %w", c.Name, err)
			}
			c.max = n
		case tagHasPrefix(m, FlagLength):
			n, err := tagInt(m)
			if err != nil {
				return nil, Errorf("when parsing length, column=%s: %w", c.Name, err)
			}
			c.length = n
		case tagHasPrefix(m, FlagEnum):
			valueSet, err := tagListContent(m)
			if err != nil {
				return nil, Errorf("column=%s: %w", c.Name, err)
			}

			c.valueSet = valueSet
			c.format = stringKindEnum
		case tagHasPrefix(m, FlagCharSet):
			valueSet, err := tagListContent(m)
			if err != nil {
				return nil, Errorf("column=%s: %w", c.Name, err)
			}
			r := []rune{}

			for _, s := range valueSet {
				if len(s) != 1 {
					return nil, Errorf("char must be of length 1, column=%s: %w", c.Name, ErrFlagFormat)
				}
				r = append(r, rune(s[0]))
			}
			c.charSet = r
		case tagHasPrefix(m, FlagTime):
			fields := tagFields(m)
			if len(fields) != 2 {
				return nil, Errorf("when parsing column=%s, tag=%s: %w", c.Name, m, ErrTagFieldNumber)
			}

			sk := stringKind(fields[1])
			if !sk.IsValid() {
				return nil, Errorf("unknown time format %s", fields[1])
			}
			c.format = sk
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

						if isTableFlag(item) {
							t.Flags = append(t.Flags, TableFlag(item))
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
					return Errorf("fk=%s: %w", c.fk, err)
				}

				tt, ok := tmap[table]
				if !ok {
					return Errorf("table=%s: %w", table, ErrFKTableNotFoundInModel)
				}

				col, err := tt.GetColumn(column)
				if err != nil {
					return Errorf("column=%s, table=%s: %w", column, table, err)
				}

				c.ForeignKey = col
			}
		}
	}
	return nil
}

func (c *Column) inferTestValuer() error {
	switch {
	case IsOneOfTypes(c.Type, StringTypeJSON):
		v, err := newValuerString(c.length, stringKindJSON, c.charSet, c.valueSet)
		if err != nil {
			return err
		}

		c.TestValuer = v
		return nil

	case IsOneOfTypes(c.Type, StringTypesAll):
		if c.format == "" {
			c.format = stringKindBasic
		}
		v, err := newValuerString(c.length, c.format, c.charSet, c.valueSet)
		if err != nil {
			return err
		}

		c.TestValuer = v
		return nil

	case IsOneOfTypes(c.Type, NumericIntegerTypesAll):
		v, err := newValuerNumeric(c.min, c.max, false)

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil

	case IsOneOfTypes(c.Type, NumericFloatTypesAll):
		v, err := newValuerNumeric(c.min, c.max, true)

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil

	case IsOneOfTypes(c.Type, BooleanTypesAll):
		v, err := newValuerBoolean()

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil

	case IsOneOfTypes(c.Type, TimeTypesAll):
		v, err := newValuerTime()

		if err != nil {
			return err
		}
		c.TestValuer = v
		return nil
	}

	return Errorf("type=%s: %w", c.Type.String(), ErrUnsuportedType)
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
				continue MainLoop
			}

			if err != nil {
				return nil, Errorf("when parsing table name: %w", err)
			}

			if x.Fields != nil {
				for _, fff := range x.Fields.List {
					if fff.Tag == nil {
						return nil, Errorf("table=%s: %w", table.Name, ErrNoColumnTag)
					}

					column, err := NewColumn(fff.Tag.Value)
					if err != nil {
						return nil, Errorf("table=%s: %w", table.Name, err)
					}
					column.Table = &table
					column.Type = info.TypeOf(fff.Type)
					column.FieldName = fff.Names[0].Name
					table.Columns = append(table.Columns, column)

					err = column.inferTestValuer()
					if err != nil {
						return nil, Errorf("when inferring test valuer - table=%s, column=%s: %w", table.StructName, column.FieldName, err)
					}
				}
			}
		}

		dbModel.Tables = append(dbModel.Tables, &table)
	}

	err = dbModel.ReconcileRelationships()
	if err != nil {
		return nil, Errorf("when reconciling relationships: %w", err)
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
