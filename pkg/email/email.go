package email

import (
	"fmt"
	"log"
	"net/smtp"
)

// can add more here as
type Email struct {
	SenderAddress string
	RecAddress    []string
	Server        string
	Port          string
	UseTTL        bool
}

func (email *Email) SetupEmailServer() {
}

func (email *Email) ExampleSendEmail() {

	connStr := email.Server + ":" + email.Port
	conn, err := smtp.Dial(connStr)
	if err != nil {
		log.Fatal(err) // TODO this probably doesn't have to be a fatal error, just log unable to setup email connection
	}

	// setup TLS
	// err = conn.StartTLS()
	if err != nil {
		fmt.Println("did not start TLS")
		log.Fatal(err)
	}

	// send the sender
	// if err := conn.Mail(email.SenderAddress); err != nil {
	// log.Fatal(err)
	// }

	// if err := conn.Rcpt(email.RecAddress[0]); err != nil {
	// 	log.Fatal(err)
	// }

	// // set the email body
	// wc, err := conn.Data()
	// if err != nil {
	// 	log.Fatal(err) // TODO this probably doesn't have to be a fatal error, just log unable to setup email connection
	// }

	// _, err = fmt.Fprintf(wc, "This is an email body")
	// if err != nil {
	// 	log.Fatal(err) // TODO this probably doesn't have to be a fatal error, just log unable to setup email connection
	// }

	// err = wc.Close()
	// if err != nil {
	// 	log.Fatal(err) // TODO this probably doesn't have to be a fatal error, just log unable to setup email connection
	// }

	// // Send the QUIT command and close the connection.
	// err = conn.Quit()
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

// func SetupEmail(email Email) {
// 	// sets up connection to email account
// 	fmt.Println("inside setup email")
// }

// func (email Email) SendEmail() {
// fmt.Println("Sending email...")
// }
