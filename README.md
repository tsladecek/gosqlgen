# gosqlgen

> [!WARNING]
> Under active development. API might change

SQL method go code generator based on table annotations using field tags.

Supported tags:
- `pk` - primary key
- `pk ai` - primary key auto incremented. Useful for inserts
- `bk` - business key
- `fk <table>.<column>` - foreign key referencing `column` on `table`
- `sd` - soft delete column. If present, the generated `delete` method will be soft delete update

## Example

Given following table spec:

```go
//go:generate go run ../../cmd/main.go -driver gosqldriver_mysql -output generatedMethods.go -outputTest generatedMethods_test.go
package dbrepo


// gosqlgen: users
type User struct {
	RawId int    `gosqlgen:"_id,pk ai"`
	Id    string `gosqlgen:"id,bk"`
	Name  string `gosqlgen:"name"`
}

// gosqlgen: addresses
type Address struct {
	RawId     int32        `gosqlgen:"_id,pk ai"`
	Id        string       `gosqlgen:"id,bk"`
	Address   string       `gosqlgen:"address,bk"`
	UserId    int          `gosqlgen:"user_id,fk users._id"`
	CountryId int          `gosqlgen:"country_id, fk countries._id"`
	DeletedAt sql.NullTime `gosqlgen:"deleted_at,sd"`
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
