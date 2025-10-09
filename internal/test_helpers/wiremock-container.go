package testhelpers

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type WiremockContainer struct {
	url string
}

func NewWiremockContainer() (WiremockContainer, error) {
	ctx := context.Background()
	localStackContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "wiremock/wiremock:3.13.1",
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForListeningPort("8080"),
		},
	})

	if err != nil {
		return WiremockContainer{}, err
	}

	host, err := localStackContainer.Host(ctx)

	if err != nil {
		return WiremockContainer{}, err
	}

	port, err := localStackContainer.MappedPort(ctx, "8080/tcp")

	if err != nil {
		return WiremockContainer{}, err
	}

	return WiremockContainer{
		url: fmt.Sprintf("http://%s:%s", host, port.Port()),
	}, nil
}

func (p *WiremockContainer) Url() string {
	return p.url
}
