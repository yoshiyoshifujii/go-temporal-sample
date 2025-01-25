package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/yoshiyoshifujii/go-temporal-sample/recovery"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

func main() {
	var workflowID, input, workflowType string
	flag.StringVar(&workflowID, "w", "trip_workflow", "WorkflowID.")
	flag.StringVar(&input, "i", "{}", "Workflow input parameters.")
	flag.StringVar(&workflowType, "wt", "tripworkflow", "Workflow type (tripworkflow|recoveryworkflow).")
	flag.Parse()

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	switch workflowType {
	case "tripworkflow":
		var userState recovery.UserState
		if err := json.Unmarshal([]byte(input), &userState); err != nil {
			log.Fatalln("Unable to unmarshal workflow input parameters", err)
		}
		workflowOptions := client.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: "recovery",
		}
		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, recovery.TripWorkflow, userState)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	case "recoveryworkflow":
		var params recovery.Params
		if err := json.Unmarshal([]byte(input), &params); err != nil {
			log.Fatalln("Unable to unmarshal workflow input parameters", err)
		}

		workflowOptions := client.StartWorkflowOptions{
			ID:                       workflowID,
			TaskQueue:                "recovery",
			WorkflowExecutionTimeout: 1 * time.Minute,
		}
		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, recovery.RecoverWorkflow, params)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	default:
		flag.PrintDefaults()
		return
	}

}
