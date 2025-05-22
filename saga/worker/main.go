package main

import (
	"log"

	"github.com/yoshiyoshifujii/go-temporal-sample/saga"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "saga", worker.Options{})
	w.RegisterWorkflow(saga.Workflow)
	w.RegisterActivity(saga.Activity01)
	w.RegisterActivity(saga.Activity02)
	w.RegisterActivity(saga.Activity03)
	w.RegisterActivity(saga.Activity01Compensation)
	w.RegisterActivity(saga.Activity02Compensation)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("failed to start worker: %v", err)
	}
}
