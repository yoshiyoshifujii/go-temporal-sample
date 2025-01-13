package main

import (
	"github.com/yoshiyoshifujii/go-temporal-sample/dsl"
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

	w := worker.New(c, "dsl", worker.Options{})

	w.RegisterWorkflow(dsl.SimpleDSLWorkflow)
	w.RegisterActivity(&dsl.SampleActivities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
