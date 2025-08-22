# gosqlgen

> [!WARNING]
> Under active development. API might change

SQL method go code generator based on table annotations using field tags.

## Table and Column definition
### Table
Table definition is expected as a comment of the struct type in following format:

`gosqlgen:table_name: REQUIRED[FLAGS]`

Where FLAGS are semicolon separated modifiers. Supported are:
- `skip tests` - tests will be skipped for this table

### Column
Column definition is expected as a field tag (similar to json tag) in following format:

`gosqlgen:"column_name: REQUIRED;sql_type:REQUIRED[FLAGS]"`
Where FLAGS are semicolon separated modifiers. Supported are:
- `pk` - primary key
- `pk ai` - primary key auto incremented. Useful for inserts
- `bk` - business key
- `fk <table>.<column>` - foreign key referencing `column` on `table`
- `sd` - soft delete column. If present, the generated `delete` method will be soft delete update

## Example

Given following table spec:

```go
//go:generate go run cmd/main.go -driver gosqldriver_mysql -out generatedMethods.go -outTest generatedMethods_test.go
package dbrepo


// gosqlgen: users
type User struct {
	RawId int    `gosqlgen:"_id;int;pk ai"`
	Id    string `gosqlgen:"id;varchar(255);bk"`
	Name  string `gosqlgen:"name;varchar(255)"`
}

// gosqlgen: addresses
type Address struct {
	RawId     int32        `gosqlgen:"_id;int;pk ai"`
	Id        string       `gosqlgen:"id;int;varchar(255);bk"`
	Address   string       `gosqlgen:"address;varchar(255);bk"`
	UserId    int          `gosqlgen:"user_id;int;fk users._id"`
	DeletedAt sql.NullTime `gosqlgen:"deleted_at;datetime;sd"`
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
