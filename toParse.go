//go:generate go run cmd/main.go -debug
package gosqlgen

import "database/sql"

// user is an object
// gosqlgen: users
type User struct {
	RawId int    `gosqlgen:"_id,pk ai"`
	Id    string `gosqlgen:"id,bk"`
	Name  int    `gosqlgen:"name"`
}

// gosqlgen: addresses
type Address struct {
	RawId     int          `gosqlgen:"_id,pk"`
	Id        string       `gosqlgen:"id,bk"`
	Address   string       `gosqlgen:"address,bk"`
	UserId    int          `gosqlgen:"user_id,fk users.id"`
	DeletedAt sql.NullTime `gosqlgen:"deleted_at,sd"`
}
