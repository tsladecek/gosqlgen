# gosqlgen: Go SQL Method Code Generator

A tool that automatically generates basic SQL methods (get, insert, update, and delete) for SQL tables declared as Go structs.
It uses Go struct comments and field tags to define the corresponding database table and columns, streamlining the creation of boilerplate database code.
Along with generated test code, this saves time writing boilerplate code/tests.

## Table and Column definition
### Table
Table definition is expected as a comment of the struct type in following format:

`gosqlgen:table_name: REQUIRED[FLAGS]`

Where FLAGS are **semicolon** separated modifiers. Supported are:
- `ignore` - methods and tests will not be generated
- `ignore update` - update method and tests for update method will not be generated
- `ignore delete` - delete method and tests for delete method woll not be generated
- `ignore test` - tests will not be generated
- `ignore test update` - tests for update method will not be generated
- `ignore test delete` - tests for delete method will not be generated

### Column
Column definition is expected as a field tag (similar to json tag) in following format:

`gosqlgen:"column_name: REQUIRED[FLAGS]"`
Where FLAGS are **semicolon** separated modifiers. Supported are:

**Column constraint flags**
- `pk` - primary key
- `ai` - auto incremented. Useful only in combination with `pk` for inserts
- `bk` - business key (<=> UNIQUE constraint)
- `fk <table>.<column>` - foreign key referencing `column` on `table`
- `sd` - soft delete column. If present, the generated `delete` method will be soft delete update

**Column value flags** - useful only for generating tests
- `min` - minimum value (relevant for numeric columns)
- `max` - maximum value (relevant for numeric columns)
- `length` - maximum length (relevant for string columns)
- `charSet{?sep} (a, b, c, d)` - alphabet (relevant for string columns). If no separator is specified (e.g. "charSet| (a | b | c | d)"), comma is used
- `enum{?sep} (val1, val2, val3)` - set of allowed values (relevant for string columns). *Format specifier*. If no separator is specified (e.g. "enum: (val1: val2: val3)"), comma is used
- `json` - string will be formatted as json (relevant for string columns). *Format specifier*
- `uuid` - string will be formatted as uuid (relevant for string columns). *Format specifier*
- `ipv4` - string will be formatted as ipv4 (xxx.xxx.xxx.xxx)
- `ipv6` - string will be formatted as ipv6 (xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx)
- `time <format>` - string will be formatted as time in given format. The format should be a string of valid go format (e.g. RFC3339, Kitchen, etc.)

> [!NOTE]
> *Format specifiers* dictate some specific format of the output strings. Only one must be supplied.
> The tool will not raise any errors if more are used within the tag. In such case, the last (right
> most) format specifier will be used

## Install

### Preferred

Download the binary from the [GitHub Releases](https://github.com/tsladecek/gosqlgen/releases) page and place it on your path

### Alternative

> [!NOTE]
> The `-version` flag wont work in this case and will only print "dev"

```shell
go install github.com/tsladecek/gosqlgen/cmd/gosqlgen@latest
```

## Example

Given following table spec:

```go
//go:generate gosqlgen -driver gosqldriver_mysql -out generatedMethods.go -outTest generatedMethods_test.go
package dbrepo

type Continent string

const (
	ContinentEurope = "Europe"
	ContinentAsia = "Asia"
	ContinentAfrica = "Africa"
)


// gosqlgen: users
type User struct {
	RawId int    `gosqlgen:"_id;pk ai"`
	Id    string `gosqlgen:"id;bk"`
	Name  string `gosqlgen:"name;length 64"`
	BirthContinent Continent `gosqlgen:"birth_continent; enum (Asia, Europe, Africa)"`
}

// gosqlgen: addresses
type Address struct {
	RawId     int32        `gosqlgen:"_id;pk ai"`
	Id        string       `gosqlgen:"id;bk"`
	Address   string       `gosqlgen:"address;bk;length 128"`
	UserId    int          `gosqlgen:"user_id;fk users._id"`
	DeletedAt sql.NullTime `gosqlgen:"deleted_at;sd"`
}
```

running

```shell
go generate .
```

will generate following methods in a file `generatedMethods.go`:

```go
func (t *User) delete(ctx context.Context, db dbExecutor) error
func (t *User) getByBusinessKeys(ctx context.Context, db dbExecutor, id string) error
func (t *User) getByPrimaryKeys(ctx context.Context, db dbExecutor, _id int) error
func (t *User) insert(ctx context.Context, db dbExecutor) error
func (t *User) updateByBusinessKeys(ctx context.Context, db dbExecutor) error
func (t *User) updateByPrimaryKeys(ctx context.Context, db dbExecutor) error

func (t *Address) delete(ctx context.Context, db dbExecutor) error
func (t *Address) getByBusinessKeys(ctx context.Context, db dbExecutor, id string, address string) error
func (t *Address) getByPrimaryKeys(ctx context.Context, db dbExecutor, _id int32) error
func (t *Address) insert(ctx context.Context, db dbExecutor) error
func (t *Address) updateByBusinessKeys(ctx context.Context, db dbExecutor) error
func (t *Address) updateByPrimaryKeys(ctx context.Context, db dbExecutor) error
```

and tests in `generatedMethods_test.go`. For the tests to work properly you have to setup the database and point the `testDb` var to the connection.

## How to use this package

The goal of the tool is to streamline the generation of basic sql methods
(get, create, update, delete) for tables, declared as structs.
If this is not your case, you can stop reading.

The package is meant to be used with `go generate`,
although it can be used as a cli tool (the `-in` argument allows for specifying input go file).
The tool generates non exported methods on provided struct types.
This is intentional and it is the responsibility of the user to export what is necessary.
It also allows for combining different methods of different struct types,
although this is not ideal since simple SQL JOIN statement is superior.

The tool will not report any illogical table definitions, unless they prevent the generation of the code.

