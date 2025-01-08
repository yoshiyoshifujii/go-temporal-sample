package main

import (
	"context"
	"github.com/pborman/uuid"
	"github.com/yoshiyoshifujii/go-temporal-sample/retryactivity/refactor_pattern1"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "retry_activity_" + uuid.New(),
		TaskQueue: "retry-activity",
	}

	bpw := refactor_pattern1.New()
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, bpw.ExecuteWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
