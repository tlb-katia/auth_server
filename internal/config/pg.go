package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	dbHost     = "PG_HOST"
	dbName     = "PG_DATABASE_NAME"
	dbUser     = "PG_USER"
	dbPassword = "PG_PASSWORD"
	dbPort     = "PG_PORT"
)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

func NewPGConfig() (PGConfig, error) {
	host := os.Getenv(dbHost)
	if len(host) == 0 {
		return nil, errors.New("pg host not found")
	}

	port := os.Getenv(dbPort)
	if len(host) == 0 {
		return nil, errors.New("pg port not found")
	}

	user := os.Getenv(dbUser)
	if len(host) == 0 {
		return nil, errors.New("pg user not found")
	}

	password := os.Getenv(dbPassword)
	if len(host) == 0 {
		return nil, errors.New("pg password not found")
	}

	dbname := os.Getenv(dbName)
	if len(host) == 0 {
		return nil, errors.New("pg name not found")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	return &pgConfig{
		dsn: psqlInfo,
	}, nil
}

func (pg *pgConfig) DSN() string {
	return pg.dsn
}
