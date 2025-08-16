//go:generate go run cmd/main.go -debug
package gosqlgen

import "database/sql"

// gosqlgen: users
type User struct {
	RawId int    `gosqlgen:"_id,pk ai"`
	Id    string `gosqlgen:"id,bk"`
	Name  int    `gosqlgen:"name"`
}

// gosqlgen: countries
type Country struct {
	RawId int    `gosqlgen:"_id,pk ai"`
	Id    string `gosqlgen:"id,bk"`
	Name  int    `gosqlgen:"name"`
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

// gosqlgen: addresses_book
type AddressBook struct {
	RawId     int    `gosqlgen:"_id,pk ai"`
	Id        string `gosqlgen:"id,bk"`
	AddressId int32  `gosqlgen:"address_id,fk addresses._id"`
}
