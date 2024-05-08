package bootstrap

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// DBConnect connects you to Postgresql based on Config.
func DBConnect(c *Config) (*sql.DB, error) {
	info := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)

	db, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
