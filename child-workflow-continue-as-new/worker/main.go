package main

import (
	cw "github.com/yoshiyoshifujii/go-temporal-sample/child-workflow-continue-as-new"
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

	w := worker.New(c, "child-workflow-continue-as-new", worker.Options{})

	w.RegisterWorkflow(cw.SampleParentWorkflow)
	w.RegisterWorkflow(cw.SampleChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
