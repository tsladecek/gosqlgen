package mattngosqlite3

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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

func getTestDB() (*sql.DB, func() error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	SetupTestDB(db, "init.sql")

	return db, db.Close
}
