package main

import (
	"context"
	"flag"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	var workflowID string
	flag.StringVar(&workflowID, "wid", "workflowID-to-cancel", "workflowID of the Workflow Execution to be canceled.")
	flag.Parse()

	if workflowID == "" {
		flag.PrintDefaults()
		return
	}

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	err = c.CancelWorkflow(context.Background(), workflowID, "")
	if err != nil {
		log.Fatalln("Unable to cancel Workflow Execution", err)
	}
	log.Println("Workflow Execution cancelled", "WorkflowID", workflowID)
}
