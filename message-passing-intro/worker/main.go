package main

import (
	"github.com/yoshiyoshifujii/go-temporal-sample/message-passing-intro"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	w := worker.New(c, message_passing_intro.TaskQueue, worker.Options{})

	w.RegisterWorkflow(message_passing_intro.GreetingWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker:", err)
	}
}
