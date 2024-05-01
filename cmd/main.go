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

	tlsConfig := &tls.Config{
		ServerName: "smtp.gmail.com",
	}

	// create email struct
	email := email.Email{
		Sender:   "raspberrypijon.tx@gmail.com",
		Password: emailPassword,
		Receiver: []string{"jonwakefield.mi@gmail.com"},
		Server:   "smtp.gmail.com",
		Port:     "587",
		UseTTL:   true,
	}

	smtpClient := email.SetupEmailServer(tlsConfig)
	defer smtpClient.Quit()
	// email.SendEmail(smtpClient)
	email.CheckTLSConnectionState(smtpClient, false)

	// email.ExampleSendEmail()
	// monitor.MainMonitor()
}
