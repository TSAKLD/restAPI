package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
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
func RedisConnect(addr string) (*redis.Client, error) {
	opts := redis.Options{
		Addr: addr,
	}

	client := redis.NewClient(&opts)

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func KafkaConnect(addr string, topic string) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", addr, topic, 0)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
