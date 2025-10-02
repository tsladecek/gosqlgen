# gosqlgen

> [!WARNING]
> Under active development. API might change

SQL method Go(lang) code generator based on table annotations using field tags.

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
- `pk ai` - primary key auto incremented. Useful for inserts
- `bk` - business key
- `fk <table>.<column>` - foreign key referencing `column` on `table`
- `sd` - soft delete column. If present, the generated `delete` method will be soft delete update

**Column value flags** - useful only for generating tests
- `min` - minimum value (relevant for numeric columns)
- `max` - maximum value (relevant for numeric columns)
- `length` - maximum length (relevant for string columns)
- `charSet (a, b, c, d)` - alphabet (relevant for string columns)
- `enum (val1, val2, val3)` - set of allowed values (relevant for string columns). *Format specifier*
- `json` - string will be formatted as json (relevant for string columns). *Format specifier*
- `uuid` - string will be formatted as uuid (relevant for string columns). *Format specifier*

*Format specifiers* dictate some specific format of the output strings. Only one must be supplied.
The tool will not raise any errors if more are used within the tag. In such case, the last (right
most) format specifier will be used

## Example

Given following table spec:

```go
//go:generate go run cmd/main.go -driver gosqldriver_mysql -out generatedMethods.go -outTest generatedMethods_test.go
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
