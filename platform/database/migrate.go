package database

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/gofiber/fiber/v2/log"
	migrate "github.com/golang-migrate/migrate/v4"
	mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func getRepoRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..")
}

func Migrate(db *sqlx.DB) error {
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return err
	}
	repoRoot := getRepoRoot()
	migrationsPath := filepath.Join(repoRoot, "migrations")

	// Convert to file URL format
	migrationsURL := "file://" + filepath.ToSlash(migrationsPath)

	fmt.Printf("Migrations URL: %s\n", migrationsURL)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"mysql", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	if err == migrate.ErrNoChange {
		log.Info("No changes in migrations")
	}
	return nil
}

func Down(db *sqlx.DB) error {
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return err
	}
	repoRoot := getRepoRoot()
	migrationsPath := filepath.Join(repoRoot, "migrations")

	// Convert to file URL format
	migrationsURL := "file://" + filepath.ToSlash(migrationsPath)

	fmt.Printf("Migrations URL: %s\n", migrationsURL)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"mysql", driver)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
