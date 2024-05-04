package monitor

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	errors.LogIfError(err)

	for _, container := range containers {
		fmt.Println("Container ID: ", container.Names)
	}

}

// test using the docker log functions
func GetLogs(containerId string) []byte {

	// get time from 24 hours ago
	startTime := time.Now().Add(-24 * time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // beautiful
	defer cancel()

	cli, _ := client.NewClientWithOpts(client.FromEnv)

	reader, err := cli.ContainerLogs(ctx, containerId, container.LogsOptions{ShowStdout: true,
		ShowStderr: true, Since: startTime.Format(time.RFC3339), Timestamps: true})

	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
