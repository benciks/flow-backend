package main

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"net/http"
	"os"

	"log"
)

//go:embed migrations
var migrations embed.FS

func main() {
	// Prepare data folders
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		if err := os.Mkdir("./data", 0755); err != nil {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", "./data/flow.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("httpfs", source, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}

	// modify for Down
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
