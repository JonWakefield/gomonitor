package main

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/joho/godotenv"
	"github.com/jonwakefield/gomonitor/pkg/email"
	"github.com/jonwakefield/gomonitor/pkg/errors"
	"github.com/jonwakefield/gomonitor/pkg/logging"
	"github.com/jonwakefield/gomonitor/pkg/monitor"
)

func main() {

	var logStartTime time.Duration = 24 // unit: hours

	// Keys:
	// container
	// network
	// volume
	// event
	// Values: start, stop, kill, die, unmount, restart etc.

	// setup custom events to listen for
	eventFilters := filters.NewArgs()
	eventFilters.Add("event", "start")
	eventFilters.Add("event", "stop")
	eventFilters.Add("event", "restart")

	err := godotenv.Load()
	errors.FatalOnErr(err)

	// setup logging
	logging.SetupLogger()

	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailSender := os.Getenv("EMAIL_SENDER")
	emailServer := os.Getenv("EMAIL_SERVER")
	emailPort := os.Getenv("EMAIL_PORT")

	recipients := []string{
		"jonwakefield.mi@gmail.com",
		// ...
	}

	// create email struct
	email := email.Email{
		Sender:     emailSender,
		Password:   emailPassword,
		Recipients: recipients,
		Server:     emailServer,
		Port:       emailPort,
	}

	email.SetupAuth()

	ctx := context.Background()

	dockerClient := monitor.CreateClient(ctx)
	defer dockerClient.Close()

	monitor.ListContainers(ctx, dockerClient)
	monitor.MonitorEvents(ctx, dockerClient, &eventFilters, &email, logStartTime)
}
