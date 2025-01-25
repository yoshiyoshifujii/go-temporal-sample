package main

import (
	"context"
	"github.com/yoshiyoshifujii/go-temporal-sample/recovery"
	"github.com/yoshiyoshifujii/go-temporal-sample/recovery/cache"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
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

	ctx := context.WithValue(context.Background(), recovery.TemporalClientKey, c)
	ctx = context.WithValue(ctx, recovery.WorkflowExecutionCacheKey, cache.NewLRU(10))

	w := worker.New(c, "recovery", worker.Options{
		BackgroundActivityContext: ctx,
	})

	w.RegisterWorkflowWithOptions(recovery.RecoverWorkflow, workflow.RegisterOptions{Name: "RecoverWorkflow"})
	w.RegisterWorkflowWithOptions(recovery.TripWorkflow, workflow.RegisterOptions{Name: "TripWorkflow"})
	w.RegisterActivity(recovery.ListOpenExecutions)
	w.RegisterActivity(recovery.RecoverExecutions)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
