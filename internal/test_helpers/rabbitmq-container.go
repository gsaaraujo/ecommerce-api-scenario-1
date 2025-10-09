package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RabbitmqContainer struct {
	url string
}

func NewRabbitmqContainer() (RabbitmqContainer, error) {
	ctx := context.Background()
	localStackContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "rabbitmq:4.1.4-management",
			ExposedPorts: []string{"5672/tcp", "15672/tcp"},
			WaitingFor:   wait.ForLog("Server startup complete").WithStartupTimeout(10 * time.Second),
			Env: map[string]string{
				"RABBITMQ_DEFAULT_USER": "guest",
				"RABBITMQ_DEFAULT_PASS": "guest",
			},
		},
	})

	if err != nil {
		return RabbitmqContainer{}, err
	}

	host, err := localStackContainer.Host(ctx)

	if err != nil {
		return RabbitmqContainer{}, err
	}

	port, err := localStackContainer.MappedPort(ctx, "5672/tcp")

	if err != nil {
		return RabbitmqContainer{}, err
	}

	return RabbitmqContainer{
		url: fmt.Sprintf("amqp://guest:guest@%s:%s", host, port.Port()),
	}, nil
}

func (p *RabbitmqContainer) Url() string {
	return p.url
}
