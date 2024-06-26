package bootstrap

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	HTTPPort string

	RedisAddr string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),

		HTTPPort: os.Getenv("HTTP_PORT"),

		RedisAddr: os.Getenv("REDIS_ADDR"),
	}, nil
}

func (c *Config) Validate() []error {
	var errorList []error

	if c.DBHost == "" {
		err := errors.New("invalid DB host field \n")
		errorList = append(errorList, err)
	}

	if c.DBPort == "" {
		err := errors.New("invalid DB port field \n")
		errorList = append(errorList, err)
	}

	if c.DBUser == "" {
		err := errors.New("invalid DB user field \n")
		errorList = append(errorList, err)
	}

	if c.DBPassword == "" {
		err := errors.New("invalid DB password field \n")
		errorList = append(errorList, err)
	}

	if c.DBName == "" {
		err := errors.New("invalid DB name field \n")
		errorList = append(errorList, err)
	}

	if c.HTTPPort == "" {
		err := errors.New("invalid HTTP port field \n")
		errorList = append(errorList, err)
	}

	if c.RedisAddr == "" {
		err := errors.New("invalid Redis host field \n")
		errorList = append(errorList, err)
	}

	if len(errorList) != 0 {
		return errorList
	}

	return nil
}
