package main

import (
	"github.com/jonwakefield/gomonitor/pkg/email"
)

func main() {

	// create email struct
	email := email.Email{
		SenderAddress: "raspberrypijon.tx@gmail.com",
		RecAddress:    []string{"jonwakefield.mi@gmail.com"},
		Server:        "smtp.gmail.com",
		Port:          "587",
		UseTTL:        true,
	}

	email.ExampleSendEmail()
	// monitor.MainMonitor()
}
