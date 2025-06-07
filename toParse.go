//go:generate go run gen.go
package main

// user is an object
// sql: users
type User struct {
	Id   int `sql:"id,pk"`
	Name int `sql:"name"`
}

// sql: addresses
type Address struct {
	Id      int    `sql:"id,pk"`
	Address string `sql:"address"`
	UserId  int    `sql:"user_id,fk users.id"`
}
