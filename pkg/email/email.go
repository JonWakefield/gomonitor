package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/smtp"
)

// can add more here as
type Email struct {
	Sender     string
	Password   string
	Recipients []string
	Server     string
	Port       string
	Auth       smtp.Auth
}

type Message struct {
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func (email *Email) SetupAuth() {
	email.Auth = smtp.PlainAuth("", email.Sender, email.Password, email.Server)
}

func (email *Email) SendEmail(msg *Message) {
	err := smtp.SendMail(fmt.Sprintf("%s:%s", email.Server, email.Port), email.Auth, email.Sender, email.Recipients, msg.ToBytes())
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
	}
	slog.Info("Successfully sent email")
}

func CreateMessage(body, subject string) *Message {
	return &Message{Subject: subject, Body: body, Attachments: make(map[string][]byte)}
}

func (msg *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(msg.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", msg.Subject))

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(fmt.Sprintf("\r\n%s\r\n", msg.Body))
	if withAttachments {
		for k, v := range msg.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
		buf.WriteString("--")
	}

	return buf.Bytes()

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
