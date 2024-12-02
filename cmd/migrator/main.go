package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	Postgres string = "postgres"
	SQLite          = "sqlite"
)

func main() {
	var storageType, storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storageType, "storage-type", "sqlite", "type of storage")
	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	if storageType == Postgres {
		applyMigrationsPostgres(storagePath, migrationsPath, migrationsTable)
	} else if storageType == SQLite {
		applyMigrationsSQLite(storagePath, migrationsPath, migrationsTable)
	} else {
		panic(fmt.Sprintf("unsupported storage type: %s", storageType))
	}
}

func applyMigrationsPostgres(storagePath, migrationsPath, migrationsTable string) {
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied")
}

func applyMigrationsSQLite(storagePath, migrationsPath, migrationsTable string) {
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied")
}

// Log represents the logger
type Log struct {
	verbose bool
}

// Printf prints out formatted string into a log
func (l *Log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Verbose shows if verbose print enabled
func (l *Log) Verbose() bool {
	return false
}
