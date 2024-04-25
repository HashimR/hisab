package store

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"main/config"
)

// NewDB initializes a new database connection pool.
func NewDB() (*sqlx.DB, error) {
	host := config.GetConfig().GetString("main-db-host")
	port := config.GetConfig().GetString("main-db-port")
	username := config.GetConfig().GetString("main-db-username")
	password := config.GetConfig().GetString("main-db-password")
	dbName := config.GetConfig().GetString("main-db-name")

	connectionURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbName)
	db, err := sqlx.Connect("mysql", connectionURL)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to DB")
		return nil, err
	}
	log.Info().Msg("Connected to db")
	return db, nil
}
