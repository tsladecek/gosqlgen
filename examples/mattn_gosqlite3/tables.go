//go:generate go run ../../cmd/gosqlgen/main.go -driver mattn_gosqlite3
package mattngosqlite3

import (
	"database/sql"
	"encoding/json"
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
	Id         string          `gosqlgen:"id;bk;length 255"`
	Name       []byte          `gosqlgen:"name"`
	payload    json.RawMessage `gosqlgen:"payload"`
	Age        sql.NullInt32   `gosqlgen:"age; min 0; max 130"`
	DrivesCar  sql.NullBool    `gosqlgen:"drives_car"`
	Birthday   sql.NullString  `gosqlgen:"birthday; time RFC3339"`
	Registered string          `gosqlgen:"registered; time RFC3339"`
}

// gosqlgen: admins; ignore test
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
	Continent Continent `gosqlgen:"continent; enum (Asia, Europe, Africa)"`
}

// gosqlgen: addresses
type Address struct {
	RawId     int32          `gosqlgen:"_id;pk;ai"`
	Id        string         `gosqlgen:"id;bk"`
	Address   string         `gosqlgen:"address"`
	UserId    int            `gosqlgen:"user_id;fk users._id"`
	CountryId int            `gosqlgen:"country_id;fk countries._id"`
	DeletedAt sql.NullString `gosqlgen:"deleted_at;sd;time RFC3339"`
	IPV4      string         `gosqlgen:"ipv4; ipv4"`
	IPV6      string         `gosqlgen:"ipv6; ipv6"`
}

// gosqlgen: addresses_book
type AddressBook struct {
	RawId     int    `gosqlgen:"_id;pk;ai"`
	Id        string `gosqlgen:"id;bk"`
	AddressId int32  `gosqlgen:"address_id;fk addresses._id"`
}
