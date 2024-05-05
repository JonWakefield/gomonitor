package apiexamples

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func DockerDiskUsage(ctx context.Context, client *client.Client) {

	options := types.DiskUsageOptions{}

	diskUsage, err := client.DiskUsage(ctx, options)
	if err != nil {
		panic(err)
	}

	fmt.Println(diskUsage.Containers)
	fmt.Println(diskUsage.Images)
	fmt.Println(diskUsage.LayersSize)
	fmt.Println(diskUsage.Volumes)

}
