package bootstrap

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"strconv"
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

// RedisConnect connects you to Redis based on Config.
func RedisConnect(c *Config) (*redis.Client, error) {
	db, err := strconv.Atoi(c.RedisDB)
	if err != nil {
		return nil, err
	}

	opts := redis.Options{
		Addr:     c.RedisHost,
		Password: c.RedisPassword,
		DB:       db,
	}

	client := redis.NewClient(&opts)

	return client, nil
}
