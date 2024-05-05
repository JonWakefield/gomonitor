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

func MonitorEvents(ctx context.Context, dockerClient *client.Client, options *filters.Args, e *email.Email, logStartTime time.Duration) {

	// TODO: modify `options` to only look for our "desired" container actions
	eventChan, errorChan := dockerClient.Events(ctx, types.EventsOptions{Filters: *options})

	slog.Info("Monitoring Docker Events...")
	for {
		select {
		case event := <-eventChan:
			// Handle event
			slog.Info(fmt.Sprintf("Received event: %s %s", event.Type, event.Action))
			body, subject := formatEmail(event)
			msg := email.CreateMessage(body, subject)
			logs := GetLogs(event.Actor.ID, logStartTime)
			msg.AttachFile(logs, event.Actor.ID)
			go e.SendEmail(msg)
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
