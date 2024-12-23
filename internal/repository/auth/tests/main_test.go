package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	dbc "github.com/tlb_katia/auth/internal/client/db"
	"github.com/tlb_katia/auth/internal/client/db/pg"
	_ "github.com/tlb_katia/auth/internal/client/db/pg"
	"log"
	"net/url"
	"testing"
	"time"
)

const testUser = "postgres"
const testPassword = "password"
const testHost = "localhost"
const testDbName = "test_db"

var db *sql.DB
var dbClient dbc.Client

func TestMain(m *testing.M) {
	pool, resource := initDB()
	m.Run()
	closeDB(pool, resource)
}

func initDB() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=my_user=%s", testUser),
			fmt.Sprintf("POSTGRES_PASSWORDr=%s", testPassword),
			fmt.Sprintf("POSTGRES_DBr=%s", testDbName),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true // Удаление контейнера после завершения тестов
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dsn := getPgDSN(resource.GetPort("5432/tcp"))

	pool.MaxWait = 30 * time.Second
	if err := pool.Retry(func() error {
		db, err = sql.Open("postgres", dsn.String())
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	configurePool()
	initMigrations()
	DBClient(dsn.String())

	return pool, resource
}

func initMigrations() {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("unable to create pg driver for migration: %v", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file:///home/katia/Desktop/auth/postgres/migrations",
		testDbName,
		driver)
	if err != nil {
		log.Fatalf("unable to create migration: %v", err)
	}

	defer func() {
		migrator.Close()
	}()

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Migrations are up to date. There is nothing to update: %v", err)
		}
		log.Fatalf("unable to apply migrations %v", err)
	}
}

func configurePool() {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
}

func getPgDSN(port string) *url.URL {
	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(testUser, testPassword),
		Host:   fmt.Sprintf("%s:%s", testHost, port),
		Path:   testDbName,
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	return pgURL
}

func closeDB(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func DBClient(DSN string) {
	ctx := context.Background()
	if dbClient == nil {
		client, err := pg.New(ctx, DSN)
		if err != nil {

		}
		dbClient = client
	}
}

/*
test functions start here
*/

func TestCreate(t *testing.T) {

}
