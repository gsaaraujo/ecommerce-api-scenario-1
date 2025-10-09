package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type LocalstackContainer struct {
	url string
}

func NewLocalstackContainer() (LocalstackContainer, error) {
	ctx := context.Background()

	localStackContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "localstack/localstack:4.2.0",
			ExposedPorts: []string{"4566/tcp"},
			WaitingFor:   wait.ForLog("Ready.").WithStartupTimeout(10 * time.Second),
			Env: map[string]string{
				"SERVICES": "secretsmanager",
			},
		},
	})

	if err != nil {
		return LocalstackContainer{}, err
	}

	host, err := localStackContainer.Host(ctx)

	if err != nil {
		return LocalstackContainer{}, err
	}

	port, err := localStackContainer.MappedPort(ctx, "4566/tcp")

	if err != nil {
		return LocalstackContainer{}, err
	}

	return LocalstackContainer{
		url: fmt.Sprintf("http://%s:%s", host, port),
	}, nil
}

func (p *LocalstackContainer) Url() string {
	return p.url
}
