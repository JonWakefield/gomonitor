package apiexamples

import (
	"context"
	"fmt"
	"io"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func DockerStats(ctx context.Context, client *client.Client) {

	containers, err := client.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println("Container: ", container.ID)

		stats, err := client.ContainerStats(ctx, container.ID, true)
		if err != nil {
			fmt.Println("Not nil!")
			panic(err)
		}
		defer stats.Body.Close()

		// read and print the stats
		buf := make([]byte, 4096)
		for {
			n, err := stats.Body.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Not nil!")
					panic(err)

				}
			}
			fmt.Print(string(buf[:n]))
		}

	}

}
