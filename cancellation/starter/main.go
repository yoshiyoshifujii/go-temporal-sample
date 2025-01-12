package main

import (
	"context"
	"flag"
	"github.com/yoshiyoshifujii/go-temporal-sample/cancellation"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	var workflowID string
	flag.StringVar(&workflowID, "w", "workflowID-to-cancel", "w is the workflowID of the workflow to be canceled.")
	flag.Parse()

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "cancel-activity",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cancellation.YourWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
