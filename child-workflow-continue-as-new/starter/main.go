package main

import (
	"context"
	"github.com/pborman/uuid"
	cw "github.com/yoshiyoshifujii/go-temporal-sample/child-workflow-continue-as-new"
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

	workflowID := "parent-workflow_" + uuid.New()
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "child-workflow-continue-as-new",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cw.SampleParentWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
