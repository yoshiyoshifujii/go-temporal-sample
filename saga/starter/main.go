package main

import (
	"context"
	"github.com/yoshiyoshifujii/go-temporal-sample/saga"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to connect to temporal server", err)
	}
	defer c.Close()

	patterns := map[string]struct {
		isActivity01Err             bool
		isActivity01CompensationErr bool
		isActivity02Err             bool
		isActivity02CompensationErr bool
		isActivity03Err             bool
	}{
		"success": {
			isActivity01Err:             false,
			isActivity01CompensationErr: false,
			isActivity02Err:             false,
			isActivity02CompensationErr: false,
			isActivity03Err:             false,
		},
		"activity01 fails": {
			isActivity01Err:             true,
			isActivity01CompensationErr: false,
			isActivity02Err:             false,
			isActivity02CompensationErr: false,
			isActivity03Err:             false,
		},
		"activity02 fails, compensation succeeds": {
			isActivity01Err:             false,
			isActivity01CompensationErr: false,
			isActivity02Err:             true,
			isActivity02CompensationErr: false,
			isActivity03Err:             false,
		},
		"activity03 fails, activity01 compensation fails": {
			isActivity01Err:             false,
			isActivity01CompensationErr: true,
			isActivity02Err:             false,
			isActivity02CompensationErr: false,
			isActivity03Err:             true,
		},
	}

	for name, p := range patterns {
		log.Printf("Running test case: %s\n", name)

		input := saga.NewWorkflowInput(
			p.isActivity01Err,
			p.isActivity01CompensationErr,
			p.isActivity02Err,
			p.isActivity02CompensationErr,
			p.isActivity03Err,
		)
		id := name
		options := client.StartWorkflowOptions{
			ID:        id,
			TaskQueue: "saga",
		}
		we, err := c.ExecuteWorkflow(context.Background(), options, saga.Workflow, input)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Printf("Started workflow: %s %s\n", we.GetID(), we.GetRunID())
		err = we.Get(context.Background(), nil)
		if err != nil {
			log.Printf("Workflow failed: %v\n", err)
		} else {
			log.Println("Workflow completed successfully")
		}
	}
}
