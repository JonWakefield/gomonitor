package email

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"strconv"

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
	errors.PanicOnErr(err)

	// setup TLS encryption
	err = client.StartTLS(tls)
	errors.PanicOnErr(err)

	err = client.Auth(auth)
	errors.PanicOnErr(err)

	return client
}

func (email *Email) SendEmail(client *smtp.Client, msg, subject string) {

	// send MAIL command to the server
	err := client.Mail(email.Sender)
	if errors.LogIfError(err) {
		return
	}

	// send email to all provided recipients
	for _, recipient := range email.Receiver {
		if err := client.Rcpt(recipient); err != nil {
			errors.LogIfError(err)
		}
	}

	// Send data command to server
	wc, err := client.Data()
	if errors.LogIfError(err) {
		return
	}
	defer wc.Close()

	emailFormatter := []byte("Subject: " + subject + "\r\n\r\n" + msg + "\r\n") // if we get fancy, this could end up being its own function

	n, err := wc.Write(emailFormatter)
	if errors.LogIfError(err) {
		return
	}

	slog.Debug("bytes written:", strconv.Itoa(n), "")
	slog.Info("Successfully sent email")
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
