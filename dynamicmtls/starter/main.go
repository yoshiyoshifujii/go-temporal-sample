package main

import (
	"context"
	"github.com/yoshiyoshifujii/go-temporal-sample/dynamicmtls"
	"github.com/yoshiyoshifujii/go-temporal-sample/helloworld"
	"go.temporal.io/sdk/client"
	"log"
	"os"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	clientOptions, err := dynamicmtls.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "hello_world_workflowID",
		TaskQueue: "hello-world-mtls",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworld.Workflow, "Temporal")
	if err != nil {
		log.Fatalf("Unable to execute workflow: %v", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result", result)
}
