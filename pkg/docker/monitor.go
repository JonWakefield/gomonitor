package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
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

func MonitorEvents(ctx context.Context, cli *client.Client) {

	options := types.EventsOptions{
		Since:   "",
		Until:   "",
		Filters: filters.Args{},
	}

	// TODO: modify `options` to only look for our "desired" container actions
	eventChan, errorChan := cli.Events(ctx, options)

	// function to monitor containers, called from a goroutine
	for {
		select {
		case event := <-eventChan:
			// Handle event
			if containerActions[event.Action] {
				fmt.Println(event.Action)
				// go email.SendEmail()
			}
		case err := <-errorChan:
			// Handle error
			fmt.Println("Received error:", err)
		}
	}
}
