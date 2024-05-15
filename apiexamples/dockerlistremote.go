// list containers on a remote server

package apiexamples

import (
	"context"
	"fmt"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func DockerListRemote() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("tcp://0.0.0.0:2376"), client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}

}
