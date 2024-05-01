package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
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

func (email *Email) SetupEmailServer(tls *tls.Config) *smtp.Client {

	connStr := email.Server + ":" + email.Port
	auth := smtp.PlainAuth("", email.Sender, email.Password, email.Server)

	// create our connection to the smtp server
	client, err := smtp.Dial(connStr)

	if err != nil {
		log.Fatal(err)
	}

	// setup TLS encryption
	err = client.StartTLS(tls)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Auth(auth)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// TODO probably want to pass in a body and subject etc.
func (email *Email) SendEmail(client *smtp.Client) {

	// send the sender
	if err := client.Mail(email.Sender); err != nil {
		log.Fatal(err)
	}
	if err := client.Rcpt(email.Receiver[0]); err != nil {
		log.Fatal(err)
	}

	// // set the email body
	wc, err := client.Data()
	if err != nil {
		log.Fatal(err) // TODO this probably doesn't have to be a fatal error, just log unable to setup email connection
	}

	_, err = fmt.Fprintf(wc, "This is from the improved email code base")
	if err != nil {
		log.Fatal(err) // TODO this probably doesn't have to be a fatal error, just log unable to setup email connection
	}
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
