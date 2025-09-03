//go:generate go run ../../cmd/main.go -driver gosqldriver_mysql -debug
package gosqldrivermysql

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
	RawId int    `gosqlgen:"_id;int;pk ai"`
	Id    string `gosqlgen:"id;varchar(255);bk"`
	Name  string `gosqlgen:"name;varchar(255);json;uuid;charset (a,b,c,d);valueset (e1,e2)"`
}

// // gosqlgen: admins;skip tests
// type Admin struct {
// 	RawId int    `gosqlgen:"_id;int;pk ai;fk users._id"`
// 	Name  string `gosqlgen:"name;varchar(255)"`
// }
//
// // gosqlgen: countries
// type Country struct {
// 	RawId     int       `gosqlgen:"_id;int;pk ai"`
// 	Id        string    `gosqlgen:"id;varchar(255);bk"`
// 	Name      string    `gosqlgen:"name;varchar(255)"`
// 	GPS       string    `gosqlgen:"gps;varchar(255)"`
// 	Continent Continent `gosqlgen:"continent;enum('Asia', 'Europe')"`
// }
//
// // gosqlgen: addresses
// type Address struct {
// 	RawId     int32        `gosqlgen:"_id;int;pk ai"`
// 	Id        string       `gosqlgen:"id;int;varchar(255);bk"`
// 	Address   string       `gosqlgen:"address;varchar(255);bk"`
// 	UserId    int          `gosqlgen:"user_id;int;fk users._id"`
// 	CountryId int          `gosqlgen:"country_id;int;fk countries._id"`
// 	DeletedAt sql.NullTime `gosqlgen:"deleted_at;datetime;sd"`
// }
//
// // gosqlgen: addresses_book
// type AddressBook struct {
// 	RawId     int    `gosqlgen:"_id;int;pk ai"`
// 	Id        string `gosqlgen:"id;varchar(255);bk"`
// 	AddressId int32  `gosqlgen:"address_id;int;fk addresses._id"`
// }
