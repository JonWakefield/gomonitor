package monitor

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/jonwakefield/gomonitor/pkg/email"
	"github.com/jonwakefield/gomonitor/pkg/errors"
)

type Set map[events.Action]bool

var containerActions = Set{
	"attach":     true,
	"start":      true,
	"kill":       true,
	"stop":       true,
	"disconnect": true,
	"die":        true,
}

func MainMonitor() {

	options := types.EventsOptions{
		Since:   "",
		Until:   "",
		Filters: filters.Args{},
	}

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

	// TODO: modify `options` to only look for our "desired" container actions
	eventChan, errorChan := cli.Events(ctx, options)

	// Process events and errors
	wg.Add(1)
	go monitorContainers(eventChan, errorChan)
	wg.Wait()
	fmt.Println("end of program")

}

func monitorContainers(eventChan <-chan events.Message, errorChan <-chan error) {
	// function to monitor containers, called from a goroutine
	for {
		select {
		case event := <-eventChan:
			// Handle event
			if containerActions[event.Action] {
				go email.SendEmail()
			}
		case err := <-errorChan:
			// Handle error
			fmt.Println("Received error:", err)
		}
	}
}

// func getActiveContainers() {
//
// }
