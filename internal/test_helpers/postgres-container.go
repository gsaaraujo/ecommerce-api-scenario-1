package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	url string
}

func NewPostgresContainer() (PostgresContainer, error) {
	ctx := context.Background()
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.2-alpine3.21",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10 * time.Second),
			Env: map[string]string{
				"POSTGRES_DB":       "postgres",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
			},
		},
	})

	if err != nil {
		return PostgresContainer{}, err
	}

	host, err := postgresContainer.Host(ctx)

	if err != nil {
		return PostgresContainer{}, err
	}

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")

	if err != nil {
		return PostgresContainer{}, err
	}

	return PostgresContainer{
		url: fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port.Port()),
	}, nil
}

func (p *PostgresContainer) Url() string {
	return p.url
}
