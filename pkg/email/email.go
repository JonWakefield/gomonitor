package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jonwakefield/gomonitor/pkg/errors"
)

// can add more here as
type Email struct {
	Sender   string
	Password string
	Receiver []string
	Server   string
	Port     string
	UseTTL   bool
}

func (email *Email) SetupSMTPClient(tls *tls.Config) *smtp.Client {

	connStr := email.Server + ":" + email.Port
	auth := smtp.PlainAuth("", email.Sender, email.Password, email.Server)

	// create our connection to the smtp server
	client, err := smtp.Dial(connStr)
	errors.FatalOnErr(err)

	// setup TLS encryption
	err = client.StartTLS(tls)
	errors.FatalOnErr(err)

	err = client.Auth(auth)
	errors.FatalOnErr(err)

	return client
}

func (email *Email) SendEmail(client *smtp.Client, msg, subject string) {

	// send MAIL command to the server
	err := client.Mail(email.Sender)
	errors.FatalOnErr(err)

	// send email to all provided recipients
	for _, recipient := range email.Receiver {
		if err := client.Rcpt(recipient); err != nil {
			log.Fatal(err)
		}
	}

	// Send data command to server
	wc, err := client.Data()
	if err != nil {
		log.Fatal(err) // TODO this doesn't have to be a fatal error, just log unable to setup email connection
	}

	defer wc.Close()

	emailFormatter := []byte("Subject: " + subject + "\r\n\r\n" + msg + "\r\n") // if we get fancy, this could end up being its own function

	n, err := wc.Write(emailFormatter)
	if err != nil {
		log.Fatal(err) // TODO this doesn't have to be a fatal error, just log unable to setup email connection
	}
	fmt.Println("bytes written: ", n)

	fmt.Println("Successfully sent email")
}

func (email *Email) CheckTLSConnectionState(client *smtp.Client, displayTLSInfo bool) bool {
	state, ok := client.TLSConnectionState()

	if displayTLSInfo {
		fmt.Println("Version is: ", state.Version)
		fmt.Println("Handshakecomplete is: ", state.HandshakeComplete)
		fmt.Println("CipherSuite is: ", state.CipherSuite)
		fmt.Println("Protocal is: ", state.NegotiatedProtocol)
		fmt.Println("ServerName is: ", state.ServerName)
	}
	return ok
}
