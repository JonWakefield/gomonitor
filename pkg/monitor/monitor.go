package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/jonwakefield/gomonitor/pkg/email"
)

type Set map[events.Action]bool

// NOTE: since we send an email on each "event", lets elimate "redudent" events (ex: kill, stop, die etc.)
var containerActions = Set{
	"start":   true,
	"restart": true,
	"pause":   true,
	"kill":    true,
	"delete":  true,
	"reload":  true,
}

func MonitorEvents(ctx context.Context, dockerClient *client.Client, e *email.Email) {

	options := types.EventsOptions{
		Since:   "",
		Until:   "",
		Filters: filters.Args{},
	}

	// TODO: modify `options` to only look for our "desired" container actions
	eventChan, errorChan := dockerClient.Events(ctx, options)

	for {
		select {
		case event := <-eventChan:
			// Handle event
			if containerActions[event.Action] {
				slog.Info(fmt.Sprintf("Registered Docker Event: %s", event.Action))
				body, subject := formatEmail(event)
				msg := email.CreateMessage(body, subject)
				logs := GetLogs(event.Actor.ID)
				msg.AttachFile(logs, event.Actor.ID)
				go e.SendEmail(msg)
			}
		case err := <-errorChan:
			// Handle error
			fmt.Println("Received error:", err)
		}
	}
}

func formatEmail(event events.Message) (string, string) {

	// unpack docker events
	action, actor, unixTime, type_ := event.Action, event.Actor, event.Time, event.Type
	id, image, name := actor.ID, actor.Attributes["image"], actor.Attributes["name"]

	// Convert Unix time to time.Time object
	timeObj := time.Unix(unixTime, 0)

	// Format time object as a human-readable string
	formattedTime := timeObj.Format("2006-01-02 15:04:05")

	subject := "Docker Event Registered"
	msg := fmt.Sprintf("Event: %s \r\nContainer Name: %s \r\nContainer ID: %s \r\nImage: %s \r\nTime: %s \r\nType: %s", action, name, id, image, formattedTime, type_)
	return msg, subject
}
