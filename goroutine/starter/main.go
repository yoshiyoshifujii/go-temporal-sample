package main

import (
	"context"
	"github.com/pborman/uuid"
	"github.com/yoshiyoshifujii/go-temporal-sample/goroutine"
	"go.temporal.io/sdk/client"
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

	workflowOptions := client.StartWorkflowOptions{
		ID:        "goroutine-" + uuid.New(),
		TaskQueue: "goroutine",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, goroutine.SampleGoroutineWorkflow, 5)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
