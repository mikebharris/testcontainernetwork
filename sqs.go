package testcontainernetwork

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"log"
)

type SqsDockerContainerConfig struct {
	Hostname       string
	ConfigFilePath string
}

type SqsDockerContainer struct {
	GenericDockerContainer
	Config SqsDockerContainerConfig
}

func (c *SqsDockerContainer) StartUsing(ctx context.Context, dockerNetwork *testcontainers.DockerNetwork) error {
	if c.Config.Hostname == "" {
		c.Config.Hostname = "sqs"
	}
	c.internalServicePort = 9324
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
	if c.dockerContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	}); err != nil {
		return fmt.Errorf("creating container: %w", err)
	}

	if err = c.dockerContainer.CopyFileToContainer(ctx, c.Config.ConfigFilePath, "/opt/elasticmq.conf", 365); err != nil {
		log.Fatalf("copying config file to docker container: %v", err)
	}

	if err := c.dockerContainer.Start(ctx); err != nil {
		return fmt.Errorf("starting container: %w", err)
	}
	return nil
}