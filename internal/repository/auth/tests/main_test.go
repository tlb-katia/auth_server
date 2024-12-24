package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	dbc "github.com/tlb_katia/auth/internal/client/db"
	"github.com/tlb_katia/auth/internal/client/db/pg"
	_ "github.com/tlb_katia/auth/internal/client/db/pg"
	"github.com/tlb_katia/auth/internal/model"
	"github.com/tlb_katia/auth/internal/repository"
	"github.com/tlb_katia/auth/internal/repository/auth"
	"log"
	"net/url"
	"testing"
	"time"
)

const testUser = "postgres"
const testPassword = "password"
const testHost = "localhost"
const testDbName = "test_db"

var dbClient dbc.Client
var repo repository.AuthRepository

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
			fmt.Sprintf("POSTGRES_USER=%s", testUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testPassword),
			fmt.Sprintf("POSTGRES_DB=%s", testDbName),
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
		dbClient, err = pg.New(context.Background(), dsn.String())
		if err != nil {
			return err
		}
		return dbClient.DB().Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	initMigrations(dsn.String())
	connectToRepoitory()

	return pool, resource
}

func initMigrations(dsn string) {
	dbCon, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database for migrations: %s", err)
	}

	if err := goose.Up(dbCon, "../../../../postgres/migrations/"); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Migrations are up to date. There is nothing to update: %v", err)
		}
		log.Fatalf("unable to apply migrations %v", err)
	}
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

func connectToRepoitory() {
	repo = auth.NewRepository(dbClient)
}

/*
test functions start here
*/

func TestCreate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		user      *model.UserInfo
		expectErr bool
	}{
		{
			name: "Valid user",
			user: &model.UserInfo{
				Name:  "John Doe",
				Email: "john.doe@example.com",
				Role:  1,
			},
			expectErr: false,
		},
		{
			name: "Duplicate email",
			user: &model.UserInfo{
				Name:  "Jane Doe",
				Email: "john.doe@example.com", // Повтор email
				Role:  0,
			},
			expectErr: true,
		},
		{
			name: "Missing name",
			user: &model.UserInfo{
				Name:  "",
				Email: "jane.doe@example.com",
				Role:  1,
			},
			expectErr: true,
		},
		{
			name: "Invalid role",
			user: &model.UserInfo{
				Name:  "Jack Doe",
				Email: "jack.doe@example.com",
				Role:  9999, // Несуществующий roleID
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userId, resErr := repo.Create(ctx, tt.user)
			if tt.expectErr {
				require.Error(t, resErr, "Ожидалась ошибка, но её не было")
			} else {
				require.NoError(t, resErr, "Неожиданная ошибка при создании пользователя")
				require.NotZero(t, userId, "Идентификатор пользователя должен быть больше 0")
			}
		})
	}

}
