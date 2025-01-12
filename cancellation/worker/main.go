package main

import (
	"github.com/yoshiyoshifujii/go-temporal-sample/cancellation"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "cancel-activity", worker.Options{})

	w.RegisterWorkflow(cancellation.YourWorkflow)
	w.RegisterActivity(&cancellation.Activities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
