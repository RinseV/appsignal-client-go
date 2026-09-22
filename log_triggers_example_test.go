package appsignal_test

import (
	"context"
	"fmt"
	"log"
	"os"

	appsignal "github.com/RinseV/appsignal-client-go"
)

func ExampleClient_GetAppLogTriggers() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	triggers, err := client.GetAppLogTriggers(context.Background(), "6aad24d1ba6bc351255e7cb5")
	if err != nil {
		log.Fatal(err)
	}

	for _, trigger := range triggers {
		fmt.Println(trigger.ID, trigger.Name, trigger.Query)
	}
}

func ExampleClient_CreateAppLogTrigger() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// The optional fields are pointers and slices; leave them nil to let the
	// API pick its own defaults.
	description := "Alert when the error log fills up"
	notificationOptions := appsignal.NotificationOptionNthInHour
	notificationTriggerValue := int32(10)

	trigger, err := client.CreateAppLogTrigger(context.Background(), appsignal.CreateAppLogTriggerInput{
		AppID:                    "6aad24d1ba6bc351255e7cb5",
		Name:                     "Too many errors",
		Query:                    "level:error",
		Description:              &description,
		NotificationOptions:      &notificationOptions,
		NotificationTriggerValue: &notificationTriggerValue,
		Severities:               []appsignal.LogSeverity{appsignal.SeverityError, appsignal.SeverityFatal},
		SourceIDs:                []string{"6aad24d1ba6bc351255e7cba"},
		NotifierIDs:              []string{"6aad24d1ba6bc351255e7cbd"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(trigger.ID)
}

func ExampleClient_UpdateAppLogTrigger() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	// Only the fields that are set get sent, so this widens the query and
	// leaves the name, severities and notifiers alone.
	query := "level:error OR level:fatal"

	trigger, err := client.UpdateAppLogTrigger(context.Background(), appsignal.UpdateAppLogTriggerInput{
		AppID:        "6aad24d1ba6bc351255e7cb5",
		LogTriggerID: "6aad24d1ba6bc351255e7cbc",
		Query:        &query,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(trigger.Query)
}

func ExampleClient_DeleteAppLogTrigger() {
	client := appsignal.NewClient("", os.Getenv("APPSIGNAL_API_TOKEN"))

	trigger, err := client.DeleteAppLogTrigger(context.Background(), appsignal.DeleteAppLogTriggerInput{
		AppID:        "6aad24d1ba6bc351255e7cb5",
		LogTriggerID: "6aad24d1ba6bc351255e7cbc",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted", trigger.Name)
}
