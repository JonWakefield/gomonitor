package main

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonwakefield/gomonitor/pkg/docker"
	"github.com/jonwakefield/gomonitor/pkg/email"
	"github.com/jonwakefield/gomonitor/pkg/errors"
	"github.com/jonwakefield/gomonitor/pkg/logging"
)

func main() {

	err := godotenv.Load()
	errors.FatalOnErr(err)

	// setup logging
	logging.SetupLogger()

	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailSender := os.Getenv("EMAIL_SENDER")

	recipients := []string{
		"jonwakefield.mi@gmail.com",
		// "buildincircuits@gmail.com",
	}

	tlsConfig := &tls.Config{
		ServerName: "smtp.gmail.com",
	}

	// create email struct
	email := email.Email{
		Sender:   emailSender,
		Password: emailPassword,
		Receiver: recipients,
		Server:   "smtp.gmail.com",
		Port:     "587",
		UseTTL:   true,
	}

	// TODO I could probably implement a recover statement if smtpClient connection fails, such that containers can still be monitored
	smtpClient := email.SetupSMTPClient(tlsConfig)
	defer smtpClient.Quit()

	// email.CheckTLSConnectionState(smtpClient, false)

	// setup connection to Docker daemon
	ctx := context.Background()

	dockerClient := docker.CreateClient(ctx)
	defer dockerClient.Close() // defer connection close until return of parent function

	docker.ListContainers(ctx, dockerClient)
	// TODO: look into a better way to send emails from the `Monitor Events` function,
	// TODO: maybe implement a channels to pass a signal saying that we need to send an email (idk)
	docker.MonitorEvents(ctx, dockerClient, smtpClient, &email)

}
