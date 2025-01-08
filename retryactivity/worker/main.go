package main

import (
	"github.com/yoshiyoshifujii/go-temporal-sample/retryactivity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	// The client and worker are heavyweight objects that should be created once and reused.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "retry-activity", worker.Options{})

	w.RegisterWorkflow(retryactivity.RetryWorkflow)
	w.RegisterActivity(retryactivity.BatchProcessingActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
