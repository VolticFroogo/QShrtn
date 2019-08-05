package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/VolticFroogo/QShrtn/config"
	"github.com/gchaincl/dotsql"
	_ "github.com/go-sql-driver/mysql" // MySQL Driver.
)

const (
	configDirectory = "configs/db.ini"
)

var (
	// SQL is the global database DB connection.
	SQL *sql.DB

	// Dot is all of the loaded queries.
	Dot *dotsql.DotSql
)

// Config is the config structure.
type Config struct {
	Name, Password, IP, Port, Database, QueriesDirectory string
}

// Init initialises the database.
func Init() (err error) {
	// Load the config.
	cfg := Config{}
	err = config.Load(configDirectory, &cfg)
	if err != nil {
		return
	}

	// Log that we are connecting to the database.
	log.Print("Connecting to database.")

	// Create the connection string from the config.
	connection := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", cfg.Name, cfg.Password, cfg.IP, cfg.Port, cfg.Database)

	// Open the SQL connection.
	SQL, err = sql.Open("mysql", connection)
	if err != nil {
		return
	}

	// Load the queries SQL file.
	Dot, err = dotsql.LoadFromFile(cfg.QueriesDirectory)
	return
}
