package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	url string
}

func NewRedisContainer() (RedisContainer, error) {
	ctx := context.Background()
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:8.2.1",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections tcp").WithStartupTimeout(10 * time.Second),
		},
	})

	if err != nil {
		return RedisContainer{}, err
	}

	host, err := postgresContainer.Host(ctx)

	if err != nil {
		return RedisContainer{}, err
	}

	port, err := postgresContainer.MappedPort(ctx, "6379/tcp")

	if err != nil {
		return RedisContainer{}, err
	}

	return RedisContainer{
		url: fmt.Sprintf("redis://:%s@%s:%s/0", "", host, port.Port()),
	}, nil
}

func (p *RedisContainer) Url() string {
	return p.url
}
