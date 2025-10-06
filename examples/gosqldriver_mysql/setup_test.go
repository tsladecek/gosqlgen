package gosqldrivermysql

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DSN struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func (d DSN) encode() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", d.User, d.Password, d.Host, d.Port, d.DBName)
}

func (d DSN) DB() *sql.DB {
	pool, err := sql.Open("mysql", d.encode())
	if err != nil {
		panic(err)
	}
	pool.SetConnMaxLifetime(time.Minute * 10)
	pool.SetMaxIdleConns(10)
	pool.SetMaxOpenConns(40)

	err = pool.Ping()

	if err != nil {
		panic(err)
	}
	return pool
}

func CreateContainer() (func(ctx context.Context), *sql.DB, error) {
	dbName := rand.Text()
	dbUser := "root"
	dbPassword := "password"
	host := "localhost"
	mappedPort := 33006
	cleanup := func(ctx context.Context) {}

	var env = map[string]string{
		"MARIADB_ROOT_PASSWORD": dbPassword,
		"MARIADB_DATABASE":      dbName,
	}
	port := "3306/tcp"

	req := testcontainers.ContainerRequest{
		Image:        "mariadb:11.7.2",
		ExposedPorts: []string{port},
		Env:          env,
		Name:         "gosqlgen_gosqldrivermysql",
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", func(host string, port nat.Port) string {
			d := DSN{User: dbUser, Password: dbPassword, Host: host, Port: port.Int(), DBName: dbName}
			return d.encode()
		}),
	}

	ctx := context.Background()
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)

	host, err = container.Host(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	mp, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	mappedPort = mp.Int()

	cleanup = func(ctx context.Context) {
		container.Terminate(ctx)
	}

	d := DSN{User: dbUser, Password: dbPassword, Host: host, Port: mappedPort, DBName: dbName}
	db := d.DB()

	return cleanup, db, nil
}

func SetupTestDB(db *sql.DB, initSQL string) {
	contentRaw, err := os.ReadFile(initSQL)
	if err != nil {
		return
	}

	var content bytes.Buffer
	for line := range bytes.SplitSeq(contentRaw, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, []byte("--")) {
			continue
		}

		content.Write(line)
		content.WriteByte('\n')
	}

	for stmt := range strings.SplitSeq(content.String(), ";") {
		stmt = strings.ReplaceAll(stmt, "\n", "")
		if stmt == "" || stmt == "\n" {
			continue
		}
		_, err = db.Exec(stmt)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var database *sql.DB

func getTestDB() (*sql.DB, func() error) {
	return database, func() error { return nil }
}

func TestMain(m *testing.M) {
	cleanup, db, err := CreateContainer()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database = db
	initSQL := "init.sql"

	SetupTestDB(db, initSQL)
	code := m.Run()
	cleanup(context.Background())
	os.Exit(code)
}
