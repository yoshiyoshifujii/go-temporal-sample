package main

import (
	"github.com/yoshiyoshifujii/go-temporal-sample/helloworldmtls"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"os"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	clientOptions, err := helloworldmtls.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world-mtls", worker.Options{})
	w.RegisterWorkflow(helloworldmtls.Workflow)
	w.RegisterActivity(helloworldmtls.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
