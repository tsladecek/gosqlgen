//go:generate go run ../../cmd/main.go -driver gosqldriver_mysql -debug
package gosqldrivermysql

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Continent string

var (
	ContinentAsia   Continent = "Asia"
	ContinentEurope Continent = "Europe"
)

type ShouldBeSkipped struct {
	somefield int
}

// gosqlgen: users
type User struct {
	RawId      int             `gosqlgen:"_id;pk;ai"`
	Id         string          `gosqlgen:"id;bk;length 5"`
	Name       []byte          `gosqlgen:"name"`
	payload    json.RawMessage `gosqlgen:"payload"`
	Age        sql.NullInt32   `gosqlgen:"age; min 0; max 130"`
	DrivesCar  sql.NullBool    `gosqlgen:"drives_car"`
	Birthday   sql.NullTime    `gosqlgen:"birthday"`
	Registered time.Time       `gosqlgen:"registered"`
}

// gosqlgen: admins;skip tests
type Admin struct {
	RawId int    `gosqlgen:"_id;pk;ai;fk users._id"`
	Name  string `gosqlgen:"name; length 31"`
}

// gosqlgen: countries
type Country struct {
	RawId     int       `gosqlgen:"_id;pk;ai"`
	Id        string    `gosqlgen:"id;bk"`
	Name      string    `gosqlgen:"name"`
	GPS       string    `gosqlgen:"gps"`
	Continent Continent `gosqlgen:"continent;valueset (Asia, Europe, Africa)"`
}

// gosqlgen: addresses
type Address struct {
	RawId     int32        `gosqlgen:"_id;pk;ai"`
	Id        string       `gosqlgen:"id;bk"`
	Address   string       `gosqlgen:"address;bk"`
	UserId    int          `gosqlgen:"user_id;fk users._id"`
	CountryId int          `gosqlgen:"country_id;fk countries._id"`
	DeletedAt sql.NullTime `gosqlgen:"deleted_at;sd"`
}

// gosqlgen: addresses_book
type AddressBook struct {
	RawId     int    `gosqlgen:"_id;pk;ai"`
	Id        string `gosqlgen:"id;bk"`
	AddressId int32  `gosqlgen:"address_id;fk addresses._id"`
}
