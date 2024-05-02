package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jonwakefield/gomonitor/pkg/errors"
)

func CreateClient(ctx context.Context) *client.Client {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	errors.PanicOnErr(err)

	return cli
}

func ListContainers(ctx context.Context, dockerClient *client.Client) {
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{})
	errors.PanicOnErr(err)

	for _, container := range containers {
		fmt.Println("Container ID: ", container.Names)
	}

}
