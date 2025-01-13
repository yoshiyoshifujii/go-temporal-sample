package main

import (
	"github.com/yoshiyoshifujii/go-temporal-sample/goroutine"
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

	w := worker.New(c, "goroutine", worker.Options{})

	w.RegisterWorkflow(goroutine.SampleGoroutineWorkflow)
	w.RegisterActivity(goroutine.Step1)
	w.RegisterActivity(goroutine.Step2)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
