package database

import (
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

func Migrate() error {
	db, err := MysqlConnection()

	if err != nil {
		return err
	}

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://platform/migrations",
		"mysql", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
