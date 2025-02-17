package main

import (
	"context"
	"flag"
	"github.com/pborman/uuid"
	"github.com/yoshiyoshifujii/go-temporal-sample/dsl"
	"go.temporal.io/sdk/client"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func main() {
	var dslConfig string
	flag.StringVar(&dslConfig, "dslConfig", "dsl/workflow1.yaml", "dslConfig specify the yaml file for the dsl workflow.")
	flag.Parse()

	data, err := os.ReadFile(dslConfig)
	if err != nil {
		log.Fatalln("failed to load dsl config file", err)
	}
	var dslWorkflow dsl.Workflow
	if err := yaml.Unmarshal(data, &dslWorkflow); err != nil {
		log.Fatalln("failed to unmarshal dsl config", err)
	}

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "dsl_" + uuid.New(),
		TaskQueue: "dsl",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, dsl.SimpleDSLWorkflow, dslWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
