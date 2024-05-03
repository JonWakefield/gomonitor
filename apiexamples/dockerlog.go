package apiexamples

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jonwakefield/gomonitor/pkg/email"
)

// test using the docker log functions
func Dockerlog(e *email.Email) {

	// this function will be called when an event with the container happens,
	// Get logs over the past 24 hours

	// get time from 24 hours ago
	startTime := time.Now().Add(-24 * time.Hour)

	containerId := "84680992cd9c"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // beautiful
	defer cancel()

	cli, _ := client.NewClientWithOpts(client.FromEnv)

	reader, err := cli.ContainerLogs(ctx, containerId, container.LogsOptions{ShowStdout: true,
		ShowStderr: true, Since: startTime.Format(time.RFC3339), Timestamps: true})

	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	body := "Message body"
	subject := "subject line"

	msg := email.CreateMessage(body, subject)

	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	msg.Attachments["Filename.txt"] = data
	fmt.Println("about to send email...")
	e.SendEmail(msg)

	// number, err := io.Copy(os.Stdout, reader)
	// if err != nil && err != io.EOF {
	// log.Fatal(err)
	// }
	// fmt.Println(number)
}
