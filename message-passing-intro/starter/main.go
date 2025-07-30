package main

import (
	"context"
	"github.com/google/uuid"
	message "github.com/yoshiyoshifujii/go-temporal-sample/message-passing-intro"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	id := "message-passing-intro-workflow-ID" + uuid.New().String()
	workflowOptions := client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: message.TaskQueue,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, message.GreetingWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow:", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	getLanguagesInput := message.GetLanguagesInput{IncludeUnsupported: false}
	supportedLangResult, err := c.QueryWorkflow(context.Background(), we.GetID(), we.GetRunID(), message.GetLanguagesQuery, getLanguagesInput)
	if err != nil {
		log.Fatalln("Unable to query workflow for supported languages:", err)
	}
	var supportedLanguages []message.Language
	err = supportedLangResult.Get(&supportedLanguages)
	if err != nil {
		log.Fatalln("Failed to get supported languages:", err)
	}
	log.Println("Supported languages:", supportedLanguages)

	langResult, err := c.QueryWorkflow(context.Background(), we.GetID(), we.GetRunID(), message.GetLanguageQuery, message.GetLanguagesInput{})
	if err != nil {
		log.Fatalln("Unable to query workflow for current language:", err)
	}
	var currentLanguage message.Language
	err = langResult.Get(&currentLanguage)
	if err != nil {
		log.Fatalln("Failed to get current language:", err)
	}
	log.Println("Current language:", currentLanguage)

	//updateHandler, err := c.UpdateWorkflow(context.Background(), client.UpdateWorkflowOptions{
	//	WorkflowID:   we.GetID(),
	//	RunID:        we.GetRunID(),
	//	UpdateName:   message.SetLanguageUpdate,
	//	WaitForStage: client.WorkflowUpdateStageAccepted,
	//	Args:         []interface{}{message.Chinese},
	//})
	//if err != nil {
	//	log.Fatalln("Unable to update workflow for language change:", err)
	//}
	//var updatePreviousLanguage message.Language
	//err = updateHandler.Get(context.Background(), &updatePreviousLanguage)
	//if err != nil {
	//	log.Fatalln("Failed to get previous language from update:", err)
	//}
	//
	//langResult, err = c.QueryWorkflow(context.Background(), we.GetID(), we.GetRunID(), message.GetLanguageQuery, message.GetLanguagesInput{})
	//if err != nil {
	//	log.Fatalln("Unable to query workflow for previous language change:", err)
	//}
	//var updateCurrentLanguage message.Language
	//err = langResult.Get(&updateCurrentLanguage)
	//if err != nil {
	//	log.Fatalln("Failed to get current language after update:", err)
	//}
	//log.Printf("Language updated from %s to %s\n", updatePreviousLanguage, updateCurrentLanguage)

	//err = c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), message.ApproveSignal, message.ApproveInput{Name: ""})
	//if err != nil {
	//	log.Fatalln("Unable to signal workflow for approval:", err)
	//}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Failed to get workflow result:", err)
	}
	log.Println("Workflow result:", result)
}
