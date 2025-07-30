package message_passing_intro

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/exp/maps"
)

type (
	Language string

	ApproveInput struct {
		Name string
	}

	GetLanguagesInput struct {
		IncludeUnsupported bool
	}
)

const (
	Chinese  Language = "chinese"
	English  Language = "english"
	French   Language = "french"
	Spanish  Language = "spanish"
	Japanese Language = "japanese"

	GetLanguagesQuery = "get-languages"
	GetLanguageQuery  = "get-language"
	SetLanguageUpdate = "set-language"
	ApproveSignal     = "approve"
)

func GreetingWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)
	approveName := "unknown"
	language := English
	greeting := map[Language]string{
		English: "Hello",
		Chinese: "你好",
	}

	err := workflow.SetQueryHandler(ctx, GetLanguagesQuery, func(input GetLanguagesInput) ([]Language, error) {
		if input.IncludeUnsupported {
			return []Language{Chinese, English, French, Spanish, Japanese}, nil
		} else {
			return maps.Keys(greeting), nil
		}
	})
	if err != nil {
		return "", err
	}

	err = workflow.SetQueryHandler(ctx, GetLanguageQuery, func(input GetLanguagesInput) (Language, error) {
		return language, nil
	})
	if err != nil {
		return "", err
	}

	updateHandlerOptions := workflow.UpdateHandlerOptions{
		Validator: func(ctx workflow.Context, newLanguage Language) error {
			if _, ok := greeting[newLanguage]; !ok {
				return fmt.Errorf("%s unsupported language", newLanguage)
			}
			return nil
		},
	}
	err = workflow.SetUpdateHandlerWithOptions(ctx, SetLanguageUpdate, func(ctx workflow.Context, newLanguage Language) (Language, error) {
		var previousLanguage Language
		previousLanguage, language = language, newLanguage
		return previousLanguage, nil
	}, updateHandlerOptions)
	if err != nil {
		return "", err
	}

	// Block until the signal is received
	var approveInput ApproveInput
	channel := workflow.GetSignalChannel(ctx, ApproveSignal)
	channel.Receive(ctx, &approveInput)
	approveName = approveInput.Name
	logger.Info("Received approval", "Approver", approveName)

	return greeting[language], nil
}
