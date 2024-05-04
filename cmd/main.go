package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonwakefield/gomonitor/pkg/email"
	"github.com/jonwakefield/gomonitor/pkg/errors"
	"github.com/jonwakefield/gomonitor/pkg/logging"
	"github.com/jonwakefield/gomonitor/pkg/monitor"
)

func main() {

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

	// setup connection to Docker daemon
	ctx := context.Background()

	dockerClient := monitor.CreateClient(ctx)
	defer dockerClient.Close() // defer connection close until return of parent function

	monitor.ListContainers(ctx, dockerClient)
	monitor.MonitorEvents(ctx, dockerClient, &email)

}
