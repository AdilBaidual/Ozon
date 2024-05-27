package containers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	Container testcontainers.Container
}

// NewRedisContainer creates and starts a new Redis container
func NewRedisContainer(ctx context.Context) (*RedisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &RedisContainer{
		Container: container,
	}, nil
}

func (r *RedisContainer) ConnectionString(ctx context.Context) (string, error) {
	host, err := r.Container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := r.Container.MappedPort(ctx, "6379")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", host, port.Port()), nil
}

func (r *RedisContainer) Terminate(ctx context.Context) error {
	return r.Container.Terminate(ctx)
}
