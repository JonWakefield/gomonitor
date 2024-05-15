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
func Dockerlog(containerId string, e *email.Email) {

	// get time from 24 hours ago
	startTime := time.Now().Add(-24 * time.Hour)

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
	fileName := containerId + "-logs" + ".txt"
	msg.Attachments[fileName] = data
	fmt.Println("about to send email...")
	e.SendEmail(msg)

}
