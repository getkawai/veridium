package docker

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	dockerclient "github.com/docker/docker/client"
	mobyclient "github.com/moby/moby/client"
)

type Client interface {
	ContainerExists(ctx context.Context, container string) (bool, container.InspectResponse, error)
	RemoveContainer(ctx context.Context, containerID string, force bool) error
	StartContainer(ctx context.Context, containerID string, containerConfig container.Config, hostConfig container.HostConfig, networkingConfig network.NetworkingConfig) error
	StopContainer(ctx context.Context, containerID string, timeout int) error
	FindContainerByLabel(ctx context.Context, label string) (string, error)
	FindAllContainersByLabel(ctx context.Context, label string) ([]string, error)
	InspectContainer(ctx context.Context, containerID string) (container.InspectResponse, error)
	ReadLogs(ctx context.Context, containerID string, options container.LogsOptions) (io.ReadCloser, error)
	ImageExists(ctx context.Context, name string) (bool, error)
	InspectImage(ctx context.Context, name string) (image.InspectResponse, error)
	PullImage(ctx context.Context, name string) error
	PullImages(ctx context.Context, names ...string) error
	CreateNetwork(ctx context.Context, name string, internal bool, labels map[string]string) error
	RemoveNetwork(ctx context.Context, name string) error
	ConnectNetwork(ctx context.Context, networkName, container, hostname string) error
	InspectVolume(ctx context.Context, name string) (volume.Volume, error)
	ReadSecrets(ctx context.Context, names []string, lenient bool) (map[string]string, error)
}

type dockerClient struct {
	apiClient func() dockerclient.APIClient
}

func NewClient(cli command.Cli) Client {
	return &dockerClient{
		apiClient: sync.OnceValue(func() dockerclient.APIClient {
			// cli.Client() returns moby/moby/client.APIClient
			// Both moby and docker client packages share the same underlying implementation
			// The concrete type (*client.Client) implements both interfaces
			mobyClient := cli.Client()

			// Safe type assertion with panic recovery
			// In practice, this should never panic as docker/cli always returns
			// a concrete *client.Client that implements both interfaces
			if c, ok := mobyClient.(dockerclient.APIClient); ok {
				return c
			}

			// This should never happen in normal operation
			// If it does, it indicates a breaking change in docker/cli
			panic("docker client type assertion failed: incompatible client types")
		}),
	}
}

func RunningInDockerCE(ctx context.Context, dockerCli command.Cli) (bool, error) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		return false, nil
	}

	mobyClient := dockerCli.Client()
	info, err := mobyClient.Info(ctx, mobyclient.InfoOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to ping Docker daemon: %w", err)
	}

	return info.Info.OperatingSystem != "Docker Desktop", nil
}
