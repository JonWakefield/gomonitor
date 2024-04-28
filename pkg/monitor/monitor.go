package monitor

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jonwakefield/gomonitor/pkg/errors"
)

func Monitor() {

	// options := types.EventsOptions{
	// 	Since:   "2024-04-25T00:00:00Z",
	// 	Until:   "2024-04-28T00:00:00Z",
	// 	Filters: filters.Args{},
	// }

	wg := sync.WaitGroup{}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	errors.PanicOnErr(err)

	defer cli.Close() // defer connection close until return of parent function

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	errors.PanicOnErr(err)

	for _, container := range containers {
		fmt.Println("Container ID: ", container.Status)
	}

	eventChan, errorChan := cli.Events(ctx, types.EventsOptions{})

	// Process events and errors
	wg.Add(1)
	go func() {
		for {
			select {
			case event := <-eventChan:
				// Handle event
				fmt.Println("Received event:", event)
			case err := <-errorChan:
				// Handle error
				fmt.Println("Received error:", err)
			}
		}
	}()
	wg.Wait()
	fmt.Println("end of program")

}
