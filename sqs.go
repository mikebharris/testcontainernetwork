package testcontainernetwork

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
)

type SqsDockerContainerConfig struct {
	Hostname   string
	Port       int
	ConfigFile string
}

type SqsDockerContainer struct {
	DockerContainer
	Config SqsDockerContainerConfig
}

func (c *SqsDockerContainer) StartUsing(ctx context.Context, dockerNetwork *testcontainers.DockerNetwork) error {
	c.internalServicePort = c.Config.Port
	req := testcontainers.ContainerRequest{
		Image:        "softwaremill/elasticmq",
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", c.internalServicePort)},
		Name:         c.Config.Hostname,
		Hostname:     c.Config.Hostname,
		Networks:     []string{dockerNetwork.Name},
		Mounts:       testcontainers.Mounts(),
		HostConfigModifier: func(config *container.HostConfig) {
			config.NetworkMode = container.NetworkMode(dockerNetwork.Name)
		},
	}
	var err error
	if c.testContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	}); err != nil {
		return fmt.Errorf("creating container: %w", err)
	}

	if err = c.testContainer.CopyFileToContainer(ctx, c.Config.ConfigFile, "/opt/elasticmq.conf", 365); err != nil {
		return fmt.Errorf("copying config file to docker container: %w", err)
	}

	if err := c.testContainer.Start(ctx); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}
	return nil
}
