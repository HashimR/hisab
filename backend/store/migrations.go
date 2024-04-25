package store

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
	"main/config"
)

func MigrateDB(migrationsPath string) error {
	host := config.GetConfig().GetString("main-db-host")
	port := config.GetConfig().GetString("main-db-port")
	username := config.GetConfig().GetString("main-db-username")
	password := config.GetConfig().GetString("main-db-password")
	dbName := config.GetConfig().GetString("main-db-name")

	connectionURL := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)
	m, err := migrate.New(
		"file://"+migrationsPath,
		connectionURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
