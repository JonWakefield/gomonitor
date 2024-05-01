package main

import (
	"crypto/tls"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonwakefield/gomonitor/pkg/email"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	emailPassword := os.Getenv("EMAIL_PASSWORD")

	recipients := []string{
		"jonwakefield.mi@gmail.com",
		// "buildincircuits@gmail.com",
	}

	tlsConfig := &tls.Config{
		ServerName: "smtp.gmail.com",
	}

	// create email struct
	email := email.Email{
		Sender:   "raspberrypijon.tx@gmail.com",
		Password: emailPassword,
		Receiver: recipients,
		Server:   "smtp.gmail.com",
		Port:     "587",
		UseTTL:   true,
	}

	smtpClient := email.SetupSMTPClient(tlsConfig)
	defer smtpClient.Quit()

	msg := "Oh no! One of your containers exploded. Ah!"
	subject := "Docker container Update"
	email.SendEmail(smtpClient, msg, subject)

	// email.CheckTLSConnectionState(smtpClient, false)

	// monitor.MainMonitor()
}
